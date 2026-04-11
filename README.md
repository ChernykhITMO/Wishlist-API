# Wishlist API

REST API для создания и управления вишлистами.

## О проекте

Пользователь может:
- зарегистрироваться и войти в систему
- создать вишлист к событию или празднику
- добавить в него подарки
- поделиться публичной ссылкой на список
- дать другому человеку возможность забронировать подарок без авторизации

Проект выполнен как тестовое задание на Go.

---

## Основные возможности

- регистрация и аутентификация по email и паролю
- хранение паролей в хэшированном виде
- CRUD для вишлистов пользователя
- CRUD для подарков внутри вишлиста
- публичный просмотр вишлиста по уникальному токену
- бронирование подарка без авторизации
- валидация входных данных
- корректные HTTP status codes и JSON-ответы
- автоматическое применение миграций при старте сервиса
- graceful shutdown

---

## Стек

- Go 1.25.1
- PostgreSQL
- Docker Compose
- `pgxpool`
- `goose`
- `golangci-lint`


## Быстрый запуск

```bash
cp .env.example .env
docker compose up --build
```

API после запуска будет доступен на `http://localhost:8080`.
Swagger UI после запуска будет доступен на `http://localhost:8080/swagger/index.html`.

Миграции руками прогонять не нужно. При старте приложение само дожидается, пока поднимется PostgreSQL, и затем применяет миграции.

## Переменные окружения

```env
HTTP_ADDR=:8080
HTTP_READ_TIMEOUT=5s
HTTP_WRITE_TIMEOUT=5s
HTTP_SHUTDOWN_TIMEOUT=5s
JWT_SECRET=change-me
POSTGRES_DB=postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_PORT=5432
DB_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
```

`DB_URL` нужен для локального запуска вне Docker.
В `docker-compose` он переопределяется автоматически на адрес контейнера с Postgres.

## Формат данных

- Поле `date` для создания и обновления вишлиста ожидается в формате `RFC3339`, например `2026-12-31T21:00:00+03:00`
- `POST /wishlists` возвращает и `id` вишлиста, и публичный `token`
- Для приватных операций используется `id`, для публичных ссылок и бронирования используется `token`

## Эндпоинты

### Авторизация

- `POST /register`
- `POST /login`

### Вишлисты

- `POST /wishlists`
- `GET /wishlists`
- `GET /wishlists/{id}`
- `PUT /wishlists/{id}`
- `DELETE /wishlists/{id}`

### Подарки

- `POST /wishlists/{wishlistId}/gifts`
- `GET /wishlists/{wishlistId}/gifts`
- `GET /wishlists/{wishlistId}/gifts/{giftId}`
- `PUT /wishlists/{wishlistId}/gifts/{giftId}`
- `DELETE /wishlists/{wishlistId}/gifts/{giftId}`

### Публичный доступ

- `GET /wishlists/public/{token}`
- `POST /wishlists/public/{token}/bookings`

## Примеры запросов

Регистрация:

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'
```

Логин:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'
```

После логина или регистрации API возвращает JWT. Его нужно передавать в `Authorization: Bearer <token>` для закрытых методов.

Создание вишлиста:

```bash
curl -X POST http://localhost:8080/wishlists \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Birthday","description":"Birthday wishlist","date":"2026-12-31T21:00:00+03:00"}'
```

Пример ответа:

```json
{
  "id": "9e9d4c7a-8d0e-4f90-8a31-6d3da6683f66",
  "token": "d9c2a9d1-e27f-4da7-89d2-f0e4b77e5f89"
}
```

Добавление подарка:

```bash
curl -X POST http://localhost:8080/wishlists/<WISHLIST_ID>/gifts \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Headphones","description":"Noise cancelling","link":"https://example.com/item","priority":5}'
```

Получение публичного вишлиста:

```bash
curl http://localhost:8080/wishlists/public/<PUBLIC_TOKEN>
```

Бронирование подарка:

```bash
curl -X POST http://localhost:8080/wishlists/public/<PUBLIC_TOKEN>/bookings \
  -H "Content-Type: application/json" \
  -d '{"giftId":"<GIFT_ID>"}'
```

## Примечание

Проект реализован полностью самостоятельно, без генерации кода с помощью LLM/AI-инструментов.  
Все архитектурные и кодовые решения, включая проектирование API, схему БД, бизнес-логику и обработку ошибок, приняты и реализованы вручную.
