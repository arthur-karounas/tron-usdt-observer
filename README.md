[🇺🇸 English](#english) | [🇷🇺 Русский](#русский)

<a name="english"></a>
# 📥 TRON USDT Observer (Enterprise Lite)

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go) 
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker)
![Security](https://img.shields.io/badge/Security-High-red)

The professional standard for TRC-20 asset surveillance.

This is a sanitized, open-source distribution of a high-performance monitoring system originally developed for institutional asset management under a private contract. It is designed for environments where reliability, data integrity, and strict access control are non-negotiable.

## 🧊 Enterprise Origin & Sanitization

This repository contains a lightweight version of a proprietary monitoring engine. The core logic has been extracted, audited for security, and decoupled from internal corporate APIs.

**Key differences from the proprietary version:**

* Removed integration mechanisms with the client's internal accounting systems.
* Replaced internal authentication services with a robust Telegram-based middleware.

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
* **Observability:** Prometheus, Grafana
* **Logging:** Zap

## 📊 Monitoring & Observability
The service is equipped with a professional monitoring stack to track engine health in real-time:

* **Prometheus Metrics:** Native instrumentation for tracking transaction throughput, API latency, and error rates.
* **Grafana Dashboard:** Pre-configured visual dashboards for instant health snapshots.

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

### 4. Monitoring Access
Once the system is running, you can access the monitoring interfaces:
* **Prometheus:** `http://localhost:9090`
* **Grafana:** `http://localhost:3000` (Default credentials: `admin/admin`)
* **Metrics Endpoint:** `http://localhost:8080/metrics`

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
# 📥 TRON USDT Observer (Enterprise Lite)

Профессиональный стандарт мониторинга активов TRC-20.

Данный проект представляет собой очищенную (sanitized) open-source версию высокопроизводительной системы мониторинга, изначально разработанной для институционального управления активами в рамках частного контракта. Система предназначена для условий, в которых надежность, целостность данных и строгий контроль доступа являются критически важными.

## 🧊 Происхождение и очистка кода

В данном репозитории представлена облегченная версия проприетарного движка мониторинга. Основная логика была выделена в этот модуль, прошла аудит безопасности и была отвязана от внутренних корпоративных API.

**Ключевые отличия от оригинальной версии:**
* Удалены механизмы интеграции со внутренними системами учёта заказчика.
* Внутренние сервисы аутентификации заменены на надежное Middleware-решение на базе Telegram.

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
* **Observability:** Prometheus, Grafana
* **Logging:** Zap

## 📊 Мониторинг и метрики
Сервис оснащен профессиональным стеком мониторинга для отслеживания состояния системы в реальном времени:

* **Prometheus Metrics:** Нативная интеграция для отслеживания пропускной способности (TPS), задержек API и частоты ошибок.
* **Grafana Dashboard:** Предустановленные дашборды для визуального контроля ключевых показателей.

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

### 4. Доступ к мониторингу
После запуска системы интерфейсы мониторинга доступны по следующим адресам:

* **Prometheus:** http://localhost:9090
* **Grafana:** http://localhost:3000 (Логин/пароль по умолчанию: admin/admin)
* **Эндпоинт метрик:** http://localhost:8080/metrics

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