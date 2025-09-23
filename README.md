[🇺🇸 English](#english) | [🇷🇺 Русский](#русский)

<a name="english"></a>
# 📥 TRON USDT Observer

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go) 
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker)
![Security](https://img.shields.io/badge/Security-High-red)

A high-performance service for monitoring incoming **USDT (TRC-20)** transactions on the Tron network with Telegram notifications.

> **Note:** This service was originally developed as an internal tool for monitoring corporate assets, which dictated strict requirements for reliability, access security, and data accuracy.

## 🛡️ Core Principles & Security

* **Non-Custodial Monitoring:** The service works exclusively with public addresses. Private keys are not used, eliminating any risk to assets
* **Access Control Middleware:** Bot management and notification delivery are restricted via Telegram ID filtering (Admin-only)
* **Transaction Deduplication:** Using Redis (SetNX) ensures exactly one notification is sent per transaction, preventing duplicates
* **Rate-Limit Protection:** Smart retry logic when interacting with the TronGrid API

## 🛠 Tech Stack

* **Language:** Go
* **Database:** PostgreSQL
* **Cache:** Redis
* **Infrastructure:** Docker & Docker Compose
* **Logging:** Zap

## 📦 Quick Start

### 1. Environment Setup

Create an `.env` file in the project root:

```env
BOT_TOKEN=your_telegram_bot_token
ADMIN_ID=your_telegram_id # You can find your ID via @userinfobot
TRON_API_KEY=your_trongrid_api_key
POLL_INTERVAL=15
```

> **Tip:** To get your Telegram `ADMIN_ID`, you can send a message to [@userinfobot](https://t.me/userinfobot) or any similar ID-checking bot in Telegram.

> **Important:** Without a `TRON_API_KEY`, public TronGrid nodes may return 429 (Rate Limit) errors. An API key is highly recommended for stable operation.

### 2. Run via Docker Compose

```bash
docker-compose up -d --build
```

### 3. Testing

The project is covered with unit tests. You can run them using:

```bash
go test ./...
```

## 🤖 Bot Commands (Admin Only)

* `/run` - Start scanner
* `/stop` - Stop scanner
* `/status` - View current configuration
* `/add_wallet <address>` - Add a wallet to monitor
* `/del_wallet <address>` - Remove a wallet from the system
* `/add_user <id>` - Grant notification access to another user
* `/del_user <id>` - Revoke notification access from a user

## 🏗 Architecture

* `cmd/bot` - Initialization and startup
* `internal/scanner` - Core: concurrent blockchain scanning
* `internal/storage` - Data layer (Postgres + Redis)
* `internal/bot` - Telegram interface logic

## 🌟 Contributing & Roadmap

* **Multi-token support:** Support any TRC-20 token via dynamic contracts in the DB
* **Web-Hook Integration:** Support webhooks for better performance
* **Reporting System:** `/report` command to generate CSV/PDF reports
* **Exchange Rates:** Integrate current USDT/Fiat exchange rates into notifications

## 📝 License

This project is licensed under the **MIT** License.

---

<a name="русский"></a>
# 📥 TRON USDT Observer

Высокопроизводительный сервис для мониторинга входящих транзакций **USDT (TRC-20)** в сети Tron с уведомлениями в Telegram.

> **Примечание:** Сервис изначально разрабатывался как внутренний инструмент для мониторинга корпоративных активов, что продиктовало строгие требования к надежности, безопасности доступа и точности данных.

## 🛡️ Принципы работы и Безопасность

* **Non-Custodial Monitoring:** Сервис работает исключительно с публичными адресами. Приватные ключи не используются, что исключает риск для активов
* **Access Control Middleware:** Управление ботом и получение уведомлений ограничено через фильтрацию по Telegram ID (Admin-only)
* **Transaction Deduplication:** Использование Redis (SetNX) гарантирует, что на каждую транзакцию придет ровно одно уведомление
* **Rate-Limit Protection:** Умная логика ретраев при взаимодействии с TronGrid API

## 🛠 Технологический стек

* **Language:** Go
* **Database:** PostgreSQL
* **Cache:** Redis
* **Infrastructure:** Docker & Docker Compose
* **Logging:** Zap

## 📦 Быстрый старт

### 1. Настройка окружения

Создайте файл `.env` в корне проекта:

```env
BOT_TOKEN=your_telegram_bot_token
ADMIN_ID=your_telegram_id # Узнать свой ID можно через @userinfobot
TRON_API_KEY=your_trongrid_api_key
POLL_INTERVAL=15
```

> **Подсказка:** Чтобы узнать свой `ADMIN_ID` в Telegram, просто напишите боту [@userinfobot](https://t.me/userinfobot) или воспользуйтесь любым аналогичным сервисом проверки ID.

> **Важно:** Без `TRON_API_KEY` публичные ноды TronGrid могут возвращать ошибку 429 (Rate Limit). Для стабильной работы сервиса ключ обязателен.

### 2. Запуск через Docker Compose

```bash
docker-compose up -d --build
```

### 3. Тестирование

Проект покрыт unit-тестами. Для их запуска выполните:

```bash
go test ./...
```

## 🤖 Команды управления (Admin Only)

* `/run` - Запустить сканер
* `/stop` - Остановить сканер
* `/status` - Просмотр текущей конфигурации
* `/add_wallet <address>` - Поставить кошелек на мониторинг
* `/del_wallet <address>` - Удалить кошелек из системы
* `/add_user <id>` - Разрешить доступ к уведомлениям другому сотруднику
* `/del_user <id>` - Забрать доступ к уведомлениям у сотрудника

## 🏗 Архитектура

* `cmd/bot` - Инициализация и запуск
* `internal/scanner` - Ядро: конкурентное сканирование блокчейна
* `internal/storage` - Слой данных (Postgres + Redis)
* `internal/bot` - Логика Telegram-интерфейса

## 🌟 Contributing & Roadmap

* **Multi-token support:** Поддержка любых TRC-20 через динамические контракты в БД
* **Web-Hook Integration:** Поддержка вебхуков для повышения производительности
* **Reporting System:** Команда `/report` для генерации отчетов в CSV/PDF
* **Exchange Rates:** Интеграция текущих курсов USDT/Fiat в уведомления

## 📝 Лицензия

Этот проект распространяется под лицензией **MIT**.