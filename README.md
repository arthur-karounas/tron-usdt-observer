[🇺🇸 English](#english) | [🇷🇺 Русский](#русский)

<a name="english"></a>
# 📥 TRON USDT Observer

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go) 
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker)
![Security](https://img.shields.io/badge/Security-Non--Custodial-blue)

A backend service for monitoring incoming TRC-20 (USDT) transactions on the TRON network.

This project is a simplified and sanitized version of a system originally developed as part of a contract engagement. It focuses on reliable transaction monitoring, data consistency, and real-time notifications.

## 🧊 Project Background & Sanitization

This repository contains a simplified version of the original system.

Some integrations and internal components were removed, and the project was adapted for open-source use.

## 🛡️ Core Principles & Security

* **Non-Custodial Monitoring:** The service works exclusively with public addresses. Private keys are not used, eliminating any risk to assets
* **Access Control Middleware:** Bot management and notification delivery are restricted via Telegram ID filtering (Admin-only)
* **Transaction Deduplication:** Using Redis (SetNX) ensures exactly one notification is sent per transaction, preventing duplicates
* **Rate-Limit Protection:** Retry logic when interacting with the TronGrid API

## 🛠 Tech Stack

* **Language:** Go
* **Database:** PostgreSQL
* **Cache:** Redis
* **Infrastructure:** Docker & Docker Compose
* **Observability:** Prometheus, Grafana
* **Logging:** Zap

## 📊 Monitoring & Observability

The service includes a monitoring stack to track system health in real-time:

* **Prometheus Metrics:** Native instrumentation for tracking transaction throughput, API latency, and error rates.
* **Grafana Dashboard:** Pre-configured visual dashboards for instant health snapshots.

## 🏗 Architecture Overview

The system processes transactions as follows:

1. A scanner continuously polls TRON wallets using the TronGrid API
2. Transactions are processed concurrently using a worker pool
3. Redis is used to check and prevent duplicate transaction processing
4. Valid transactions are stored in PostgreSQL
5. Notifications are sent via Telegram bot

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

## 🏗 Project Structure

* `cmd/bot` - Initialization and startup
* `internal/scanner` - Core: concurrent blockchain scanning
* `internal/storage` - Data layer (Postgres + Redis)
* `internal/bot` - Telegram interface logic
* `internal/obs` - Metrics for observability

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

Сервис на backend для мониторинга входящих TRC-20 (USDT) транзакций в сети TRON.

Данный проект представляет собой упрощённую и очищенную версию системы, разработанной в рамках контрактной работы. Основной фокус - надёжный мониторинг транзакций, целостность данных и отправка уведомлений в реальном времени.


## 🧊 Происхождение проекта и очистка

В данном репозитории представлена упрощённая версия исходной системы.

Некоторые интеграции и внутренние компоненты были удалены, а проект адаптирован для публикации в open-source.

## 🛡️ Принципы работы и Безопасность

* **Non-Custodial Monitoring:** Сервис работает только с публичными адресами. Приватные ключи не используются, что исключает риск для средств
* **Контроль доступа:** Управление ботом и получение уведомлений ограничены через фильтрацию по Telegram ID (только для администратора)
* **Дедупликация транзакций:** Использование Redis (SetNX) гарантирует, что каждая транзакция обрабатывается только один раз
* **Защита от лимитов API:** Реализована логика повторных попыток при работе с TronGrid API

## 🛠 Технологический стек

* **Язык:** Go
* **База данных:** PostgreSQL
* **Кэш:** Redis
* **Инфраструктура:** Docker и Docker Compose
* **Мониторинг:** Prometheus, Grafana
* **Логирование:** Zap

## 📊 Мониторинг и метрики
Сервис включает стек мониторинга для отслеживания состояния системы в реальном времени:

* **Метрики Prometheus:** Отслеживание пропускной способности, задержек API и количества ошибок
* **Дашборды Grafana:** Визуализация ключевых метрик системы

## 🏗 Архитектура

Система работает следующим образом:

1. Сканер регулярно опрашивает кошельки TRON через TronGrid API
2. Транзакции обрабатываются параллельно с использованием пула воркеров
3. Redis используется для проверки и предотвращения повторной обработки транзакций
4. Валидные транзакции сохраняются в PostgreSQL
5. Уведомления отправляются через Telegram-бота

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

## 🏗 Структура проекта

* `cmd/bot` - Инициализация и запуск
* `internal/scanner` - Ядро: конкурентное сканирование блокчейна
* `internal/storage` - Слой данных (Postgres + Redis)
* `internal/bot` - Логика Telegram-интерфейса
* `internal/obs` - Определение метрик для observability

## 🌟 Contributing & Roadmap

* **Multi-token support:** Поддержка любых TRC-20 через динамические контракты в БД
* **Web-Hook Integration:** Поддержка вебхуков для повышения производительности
* **Reporting System:** Команда `/report` для генерации отчетов в CSV/PDF
* **Exchange Rates:** Интеграция текущих курсов USDT/Fiat в уведомления

## 📝 Лицензия

Этот проект распространяется под лицензией **MIT**.