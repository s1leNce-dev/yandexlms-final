# 📡 Калькулятор распределённых вычислений

Распределённый вычислитель арифметических выражений. Финал

---

## 🚀 Функциональность

- Регистрация и вход пользователей
- Защищённый доступ к API через JWT (в cookie)
- Отправка выражений для вычислений (поддержка `+ - * /`)
- Разделение выражений на подзадачи
- Распределённое выполнение через gRPC-агентов
- Хранение и просмотр истории вычислений
- Интеграционные и модульные тесты

---

## 📦 Технологии

- Язык: Go
- Веб-фреймворк: Gin
- БД: SQLite + GORM
- Авторизация: JWT + Cookie
- gRPC: для взаимодействия с агентами
- Docker: для контейнеризации
- Тестирование: `go test`

---

## 🔧 Запуск проекта

### ✅ Вариант 1: Docker Compose

```bash
docker-compose up --build
````

* API будет доступен по адресу: `http://localhost:8000/api/v1`
* gRPC сервис будет запущен на порту: `:50051`

---

### 🧪 Вариант 2: Ручной запуск

#### 1. Убедитесь, что установлен Go 1.20+

#### 2. Запуск сервиса

```bash
go run ./cmd/calc_service/...
```

#### 3. Запуск gRPC агента (в отдельной вкладке терминала):

```bash
go run ./cmd/agent/...
```

---

## 🔐 Авторизация

### Регистрация

```bash
curl --location 'http://localhost:8000/api/v1/regin' \
--header 'Content-Type: application/json' \
--data '{
  "login": "user1",
  "password": "pass123"
}'
```

✅ Успешный ответ:

```json
{ "message": "success" }
```

---

### Вход

```bash
curl --location 'http://localhost:8000/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{
  "login": "user1",
  "password": "pass123"
}'
```

✅ Успешный ответ:

```json
{ "message": "success" }
```

*После авторизации в cookie устанавливаются JWT-токены.*

---

## 📡 Примеры API-запросов

### 📥 Отправка выражения

```bash
curl --location 'http://localhost:8000/api/v1/user/calculate' \
--header 'Content-Type: application/json' \
--header 'Cookie: jwt_access=<ваш_токен>' \
--data '{
  "expression": "2 + 2 * 2"
}'
```

✅ Успешный ответ:

```json
{ "id": 1 }
```

❌ Ошибка:

```json
{ "error": "Invalid expression format" }
```

---

### 📄 Получение всех выражений

```bash
curl --location 'http://localhost:8000/api/v1/user/expressions' \
--header 'Cookie: jwt_access=<ваш_токен>'
```

✅ Пример ответа:

```json
{
  "expressions": [
    { "ID": 1, "Status": "done", "Result": "6" },
    { "ID": 2, "Status": "pending", "Result": null }
  ]
}
```

---

### 🔍 Получение выражения по ID

```bash
curl --location 'http://localhost:8000/api/v1/user/expressions/1' \
--header 'Cookie: jwt_access=<ваш_токен>'
```

✅ Ответ:

```json
{
  "expression": {
    "ID": 1,
    "Status": "done",
    "Result": "6",
    "Tasks": [...]
  }
}
```

❌ Ошибка:

```json
{ "error": "Expression not found" }
```

---

## 🧪 Тестирование

Модульные и интеграционные тесты:

```bash
go test ./...
```

---

## ⚙ Переменные окружения

```env
# Server
SERVER_HOST="0.0.0.0"
SERVER_PORT="8000"

# JWT_SECRETS
JWT_ACCESS_TOKEN="jwt_access228"
JWT_REFRESH_TOKEN="jwt_refresh1337"
```

---

## 📁 Структура

```text
cmd/
  ├─ calc_service/     # Запуск REST API + gRPC-сервера
  └─ agent/            # gRPC агент

internal/
  ├─ handlers/
  ├─ middlewares/
  ├─ models/
  ├─ proto/
  └─ utils/

docker-compose.yml
```

