package bot

import (
	"time"

	"github.com/arthur-karounas/tron-usdt-observer/internal/config"
	"github.com/arthur-karounas/tron-usdt-observer/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

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

type Bot struct {
	b       *tele.Bot
	db      Store
	scanner ScannerController
	logger  *zap.SugaredLogger
	cfg     *config.Config
}

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

func (bot *Bot) Start() {
	bot.logger.Info("Telegram bot started")
	bot.b.Start()
}

func (bot *Bot) Stop() {
	bot.b.Stop()
}

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

func (bot *Bot) setupHandlers() {
}
