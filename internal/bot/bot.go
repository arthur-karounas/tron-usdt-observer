package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/arthur-karounas/tron-usdt-observer/internal/config"
	"github.com/arthur-karounas/tron-usdt-observer/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

// External dependencies interfaces
type Store interface {
	GetUsers() ([]int64, error)
	GetWallets() ([]storage.TrackedWallet, error)
	AddWallet(address string) error
	RemoveWallet(address string) error
	AddUser(id int64) error
	RemoveUser(id int64) error
}

type ScannerController interface {
	SetRunning(state bool)
	IsRunning() bool
}

// Bot represents the Telegram controller
type Bot struct {
	b       *tele.Bot
	db      Store
	scanner ScannerController
	logger  *zap.SugaredLogger
	cfg     *config.Config
}

// New creates and configures a new Bot instance
func New(cfg *config.Config, db Store, scn ScannerController, logger *zap.SugaredLogger) (*Bot, error) {
	pref := tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	botInst := &Bot{b: b, db: db, scanner: scn, logger: logger, cfg: cfg}
	botInst.setupHandlers()

	return botInst, nil
}

// Lifecycle management
func (bot *Bot) Start() {
	bot.logger.Info("Telegram bot started")
	bot.b.Start()
}

func (bot *Bot) Stop() {
	bot.b.Stop()
}

// Broadcast notifications to all authorized users
func (bot *Bot) SendNotification(msg string) {
	users, err := bot.db.GetUsers()
	if err != nil {
		bot.logger.Errorf("Failed to get users for notification: %v", err)
		return
	}

	for _, userID := range users {
		user := &tele.User{ID: userID}
		_, err := bot.b.Send(user, msg, tele.ModeHTML, tele.NoPreview)
		if err != nil {
			bot.logger.Errorf("Failed to send notification to %d: %v", userID, err)
		}
	}
}

// Middleware: restrict access to admin only
func (bot *Bot) adminOnly(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Sender().ID != bot.cfg.AdminID {
			return nil
		}
		return next(c)
	}
}

// Define command routes and middleware
func (bot *Bot) setupHandlers() {
	adminGroup := bot.b.Group()
	adminGroup.Use(bot.adminOnly)

	adminGroup.Handle("/start", bot.handleStart)
	adminGroup.Handle("/status", bot.handleStatus)
	adminGroup.Handle("/run", bot.handleRun)
	adminGroup.Handle("/stop", bot.handleStop)
	adminGroup.Handle("/add_wallet", bot.handleAddWallet)
	adminGroup.Handle("/del_wallet", bot.handleDelWallet)
	adminGroup.Handle("/add_user", bot.handleAddUser)
	adminGroup.Handle("/del_user", bot.handleDelUser)
}

// --- Command Handlers ---

func (bot *Bot) handleStart(c tele.Context) error {
	text := "System online. Scanner controls:\n" +
		"/run - Start scanner\n" +
		"/stop - Stop scanner\n" +
		"/status - View current configuration\n\n" +
		"Management:\n" +
		"/add_wallet <address>\n" +
		"/del_wallet <address>\n" +
		"/add_user <id>\n" +
		"/del_user <id>"
	return c.Send(text)
}

func (bot *Bot) handleRun(c tele.Context) error {
	if bot.scanner.IsRunning() {
		return c.Send("Scanner is already running.")
	}
	bot.scanner.SetRunning(true)
	return c.Send("Scanner started.")
}

func (bot *Bot) handleStop(c tele.Context) error {
	if !bot.scanner.IsRunning() {
		return c.Send("Scanner is already stopped.")
	}
	bot.scanner.SetRunning(false)
	return c.Send("Scanner stopped.")
}

func (bot *Bot) handleStatus(c tele.Context) error {
	status := "Stopped 🔴"
	if bot.scanner.IsRunning() {
		status = "Running 🟢"
	}

	wallets, _ := bot.db.GetWallets()
	users, _ := bot.db.GetUsers()

	msg := fmt.Sprintf("<b>Scanner Status:</b> %s\n\n<b>Tracked Wallets (%d):</b>\n", status, len(wallets))
	for _, w := range wallets {
		msg += fmt.Sprintf("• <code>...%s</code>\n", w.Address[len(w.Address)-4:])
	}

	msg += fmt.Sprintf("\n<b>Authorized Listeners (%d):</b>\n", len(users))
	for _, u := range users {
		msg += fmt.Sprintf("• <code>%d</code>\n", u)
	}

	return c.Send(msg, tele.ModeHTML)
}

// --- Wallet Management ---

func (bot *Bot) handleAddWallet(c tele.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send("Usage: /add_wallet <address>")
	}

	err := bot.db.AddWallet(args[0])
	if err != nil {
		return c.Send("Database error.")
	}
	return c.Send(fmt.Sprintf("Address added: ...%s", args[0][len(args[0])-4:]))
}

func (bot *Bot) handleDelWallet(c tele.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send("Usage: /del_wallet <address>")
	}

	err := bot.db.RemoveWallet(args[0])
	if err != nil {
		bot.logger.Errorf("Database error removing wallet: %v", err)
		return c.Send("Database error. Could not remove wallet.")
	}

	return c.Send(fmt.Sprintf("Address removed: ...%s", args[0][len(args[0])-4:]))
}

// --- User Management ---

func (bot *Bot) handleAddUser(c tele.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send("Usage: /add_user <user_id>")
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return c.Send("Error: user_id must be a number.")
	}
	bot.db.AddUser(id)
	return c.Send(fmt.Sprintf("User added: %d", id))
}

func (bot *Bot) handleDelUser(c tele.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send("Usage: /del_user <user_id>")
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return c.Send("Error: user_id must be a number.")
	}
	if id == bot.cfg.AdminID {
		return c.Send("Error: Cannot remove the main administrator.")
	}

	err = bot.db.RemoveUser(id)
	if err != nil {
		bot.logger.Errorf("Database error removing user: %v", err)
		return c.Send("Database error. Could not remove user.")
	}

	return c.Send(fmt.Sprintf("User removed: %d", id))
}
