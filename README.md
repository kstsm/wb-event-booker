# EventBooker

## Описание

Сервис предназначен для **бронирования мест на мероприятиях** с поддержкой автоматической отмены неоплаченных бронирований по истечении заданного времени.

**Основные возможности:**
- создание платных и бесплатных мероприятий с настраиваемым сроком жизни бронирования
- бронирование мест пользователями
- подтверждение бронирования (оплата)
- автоматическая отмена неоплаченных бронирований через фоновый планировщик
- уведомления пользователей об отмене бронирования через Telegram
- поддержка множественных пользователей и их регистрации
- веб-интерфейс для пользователей и администраторов

# Фоновый планировщик

Сервис автоматически обрабатывает просроченные бронирования через фоновый планировщик, который:

- запускается при старте сервиса
- периодически проверяет бронирования со статусом `reserved`, у которых истек срок (deadline)
- автоматически отменяет просроченные бронирования с использованием транзакций
- отправляет уведомления пользователям через Telegram (если настроен токен бота)

Интервал проверки настраивается через переменную окружения `SCHEDULER_CHECK_INTERVAL` (в секундах).

# Telegram уведомления

Если в `.env` файле указан `TELEGRAM_BOT_TOKEN`, сервис будет отправлять уведомления
пользователям об отмене бронирования через
Telegram Bot API. Пользователь должен иметь указанный `telegram_id` при регистрации.

### HTTP API

- POST /api/events - создание мероприятия
- POST /api/users - создание пользователя
- POST /api/events/{id}/book - бронирование места
- POST /api/events/{id}/confirm - подтверждение (оплата) брони
- GET /api/events/{id} - получение информации о мероприятии и свободных местах
- GET /api/events - получение списка всех мероприятий
- GET /api/events/{id}/bookings - получение списка бронирований мероприятия


## Установка и запуск проекта

### 1. Клонирование репозитория

```bash
git clone https://github.com/kstsm/wb-event-booker
```

### 2. Настройка переменных окружения

Создайте `.env` файл, скопировав в него значения из `.example.env`:

```bash
cp .example.env .env
```

Отредактируйте `.env` файл, указав необходимые значения:

```bash
# Server
SRV_HOST=localhost
SRV_PORT=8080

# Postgres
POSTGRES_CONTAINER_NAME=event-booking-db
POSTGRES_VOLUME_NAME=event_booking_data
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=admin
POSTGRES_PASSWORD=admin
POSTGRES_DB=event_booker
POSTGRES_SSLMODE=disable

# Telegram (опционально, для уведомлений)
TELEGRAM_BOT_TOKEN=TOKEN

# Scheduler (интервал проверки просроченных бронирований в секундах)
SCHEDULER_CHECK_INTERVAL=10

# Goose
DB_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}

MIGRATIONS_DIR=./internal/migrations
```

### 3. Запуск зависимостей (Docker)

```bash
make docker-up
```

Это запустит PostgreSQL в контейнере Docker.

### 4. Миграция базы данных

```bash
make migrate-up
```

### 5. Запуск сервиса

```bash
make run
```

Сервис будет доступен по адресу: http://localhost:8080

# API запросы

## POST /api/events - Создание мероприятия

**URL:** `http://localhost:8080/api/events`

**Content-Type:** `application/json`

**Параметры:**

- `name` (обязательно) - название мероприятия
- `date` (обязательно) - дата и время проведения в формате ISO 8601
- `total_seats` (обязательно) - общее количество мест
- `booking_lifetime_hours` (обязательно) - срок жизни бронирования в часах
- `booking_lifetime_minutes` (обязательно) - срок жизни бронирования в минутах
- `requires_payment_confirmation` (обязательно) — требуется ли подтверждение оплаты (true/false)

**Body:**

```json
{
  "name": "Концерт классической музыки",
  "date": "2025-12-15T19:00:00Z",
  "total_seats": 100,
  "booking_lifetime_hours": 2,
  "booking_lifetime_minutes": 0,
  "requires_payment_confirmation": true
}
```

**Ожидаемый ответ (201 Created):**

```json
{
  "event": {
    "id": "fcdcf25c-fbc1-4941-a3b7-40a24bb71446",
    "name": "Golang Meetup Wildberries",
    "date": "2025-12-15T19:00:00Z",
    "total_seats": 100,
    "reserved_seats": 0,
    "booked_seats": 0,
    "booking_lifetime": 120,
    "requires_payment_confirmation": true,
    "created_at": "2025-12-02T16:44:42.788089761Z"
  },
  "message": "event created successfully"
}
```

### Ошибки:

**Некорректный JSON (400 Bad Request):**

```json
{
  "error": "invalid request body"
}
```

**Ошибки валидации (400 Bad Request):**

```json
{
  "error": "event name is required"
}
```

```json
{
  "error": "total number of seats must be greater than or equal to 1"
}
```

```json
{
  "error": "booking lifetime hours cannot be negative"
}
```

```json
{
  "error": "booking lifetime minutes must be between 0 and 59"
}
```

```json
{
  "error": "minimum booking lifetime is 1 minutes"
}
```

```json
{
  "error": "booking lifetime cannot be negative"
}
```

```json
{
  "error": "invalid date format"
}
```

```json
{
  "error": "event date cannot be in the past"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---
## POST /api/users - Создание пользователя

**URL:** `http://localhost:8080/api/users`

**Content-Type:** `application/json`

**Параметры:**

- `name` (обязательно) — имя пользователя
- `email` (обязательно) — email пользователя (уникальный)
- `telegram_id` (опционально) — Telegram ID для уведомлений

**Body:**

```json
{
  "name": "Иван Иванов",
  "email": "Ivan@gmail.com",
  "telegram_id": 123456788
}
```
**Ожидаемый ответ (201 Created):**

```json
{
  "user": {
    "id": "342540df-bb18-4c4f-8c00-f17ed9045bee",
    "name": "Иван Иванов",
    "email": "Ivan@gmail.com",
    "telegram_id": 123456788,
    "created_at": "2025-12-02T22:55:01.769582756+06:00"
  },
  "message": "user created successfully"
}
```
### Ошибки:

**Некорректный JSON (400 Bad Request):**

```json
{
  "error": "invalid request body"
}
```

**Ошибки валидации (400 Bad Request):**

```json
{
  "error": "user name is required"
}
```

```json
{
  "error": "name must contain only letters"
}
```

```json
{
  "error": "user email is required"
}
```

```json
{
  "error": "invalid email format"
}
```

```json
{
  "error": "telegram id must be >= 1000000"
}
```

```json
{
  "error": "telegram id must be <= 9999999999"
}
```

**Email уже существует (400 Bad Request):**

```json
{
  "error": "email already exists"
}
```

**Telegram ID уже используется (400 Bad Request):**

```json
{
  "error": "telegram id already exists"
}
```

## POST /api/events/{id}/book - Бронирование места

**URL:** `http://localhost:8080/api/events/{id}/book`

**Content-Type:** `application/json`

**Параметры:**

- `{id}` (обязательно)   - UUID мероприятия
- `email` (обязательно) - email пользователя для бронирования

**Body:**

```json
{
  "email": "Ivan@gmail.com"
}
```

**Ожидаемый ответ (201 Created):**

```json
{
  "booking_id": "3d590492-d3ed-4afb-b6c6-566ede90c7e7",
  "deadline": "2025-12-02T19:00:50Z",
  "message": "booking created successfully"
}
```

### Ошибки:

**Некорректный ID мероприятия (400 Bad Request):**

```json
{
  "error": "invalid id"
}
```

или

```json
{
  "error": "id is required"
}
```

**Некорректный JSON (400 Bad Request):**

```json
{
  "error": "invalid request body"
}
```

**Мероприятие не найдено (404 Not Found):**

```json
{
  "error": "event not found"
}
```

**Пользователь не найден (404 Not Found):**

```json
{
  "error": "user not found"
}
```

**Нет свободных мест (409 Conflict):**

```json
{
  "error": "no available seats"
}
```

**Пользователь уже забронировал это мероприятие (409 Conflict):**

```json
{
  "error": "user already has a booking for this event"
}
```

**Мероприятие уже прошло (400 Bad Request):**

```json
{
  "error": "event has expired"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## POST /api/events/{id}/confirm - Подтверждение бронирования

**URL:** `http://localhost:8080/api/events/{id}/confirm`

**Content-Type:** `application/json`

**Параметры:**
- `{id}` (обязательно) - UUID мероприятия
- `booking_id` (обязательно) - UUID бронирования для подтверждения

**Body:**

```json
{
  "booking_id": "3d590492-d3ed-4afb-b6c6-566ede90c7e7"
}
```
**Ожидаемый ответ (200 OK):**

```json
{
  "message": "confirmed successfully"
}
```

### Ошибки:

**Некорректный ID мероприятия (400 Bad Request):**

```json
{
  "error": "invalid id"
}
```

или

```json
{
  "error": "id is required"
}
```

**Некорректный JSON (400 Bad Request):**

```json
{
  "error": "invalid request body"
}
```

**Бронирование не найдено (404 Not Found):**

```json
{
  "error": "booking not found"
}
```

**Бронирование не в статусе reserved (400 Bad Request):**

```json
{
  "error": "booking is not in reserved status"
}
```

**Срок бронирования истек (400 Bad Request):**

```json
{
  "error": "booking deadline has passed"
}
```

**Мероприятие не требует подтверждения оплаты (400 Bad Request):**

```json
{
  "error": "event does not require payment confirmation"
}
```

**Мероприятие уже прошло (400 Bad Request):**

```json
{
  "error": "event has expired"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## GET /api/events/{id} - Получение информации о мероприятии

**URL:** `http://localhost:8080/api/events/{id}`

**Поля ответа:**
- `reserved_seats` - количество зарезервированных (неоплаченных) мест
- `booked_seats` - количество подтверждённых (оплаченных) мест

**Ожидаемый ответ (200 OK):**

```json
{
  "event": {
    "id": "fcdcf25c-fbc1-4941-a3b7-40a24bb71446",
    "name": "Golang Meetup Wildberries",
    "date": "2025-12-16T01:00:00+06:00",
    "total_seats": 100,
    "reserved_seats": 0,
    "booked_seats": 0,
    "booking_lifetime": 120,
    "requires_payment_confirmation": true,
    "created_at": "2025-12-02T22:44:42.788089+06:00"
  }
}
```

### Ошибки:

**Некорректный ID мероприятия (400 Bad Request):**

```json
{
  "error": "invalid id"
}
```

или

```json
{
  "error": "id is required"
}
```

**Мероприятие не найдено (404 Not Found):**

```json
{
  "error": "event not found"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## GET /api/events - Получение списка всех мероприятий

**URL:** `http://localhost:8080/api/events`

**Ожидаемый ответ (200 OK):**

```json
{
  "events": [
    {
      "id": "3dcb4cdd-d45c-4f1c-9b3b-67063a70874b",
      "name": "Концерт классической музыки",
      "date": "2025-12-16T01:00:00+06:00",
      "total_seats": 100,
      "reserved_seats": 0,
      "booked_seats": 0,
      "booking_lifetime": 120,
      "requires_payment_confirmation": true,
      "created_at": "2025-12-02T22:38:53.719068+06:00"
    },
    {
      "id": "fcdcf25c-fbc1-4941-a3b7-40a24bb71446",
      "name": "Golang Meetup Wildberries",
      "date": "2025-12-16T01:00:00+06:00",
      "total_seats": 100,
      "reserved_seats": 0,
      "booked_seats": 0,
      "booking_lifetime": 120,
      "requires_payment_confirmation": true,
      "created_at": "2025-12-02T22:44:42.788089+06:00"
    },
    {
      "id": "b1079863-45c1-4a40-a2ae-5196b0bf1ed0",
      "name": "Golang Meetup Wildberries",
      "date": "2025-12-16T01:00:00+06:00",
      "total_seats": 100,
      "reserved_seats": 0,
      "booked_seats": 1,
      "booking_lifetime": 120,
      "requires_payment_confirmation": true,
      "created_at": "2025-12-02T22:58:50.532607+06:00"
    }
  ]
}
```

### Ошибки:

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---

## GET /api/events/{id}/bookings - Получение списка бронирований мероприятия

**URL:** `http://localhost:8080/api/events/{id}/bookings`

**Статусы бронирования:**

- `reserved` - зарезервировано (ожидает оплаты)
- `confirmed` - подтверждено (оплачено)
- `cancelled` - отменено (автоматически)

**Ожидаемый ответ (200 OK):**

```json
{
  "bookings": [
    {
      "id": "9181148a-0f66-4a55-b89c-0aeba08800e2",
      "event_id": "fcdcf25c-fbc1-4941-a3b7-40a24bb71446",
      "user_id": "342540df-bb18-4c4f-8c00-f17ed9045bee",
      "status": "confirmed",
      "deadline": "2025-12-03T01:36:07.704986+06:00",
      "created_at": "2025-12-02T23:36:07.707672+06:00",
      "updated_at": "2025-12-02T23:36:18.88273+06:00"
    }
  ]
}
```

### Ошибки:

**Некорректный ID мероприятия (400 Bad Request):**

```json
{
  "error": "invalid id"
}
```

или

```json
{
  "error": "id is required"
}
```

**Внутренняя ошибка сервера (500 Internal Server Error):**

```json
{
  "error": "internal server error"
}
```

---
