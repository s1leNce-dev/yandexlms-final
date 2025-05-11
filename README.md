Вот **чистый `README.md`** на основе нового кода, без лишних пояснений:

````markdown
# 📡 Калькулятор распределённых вычислений

Распределённая система вычислений с авторизацией, REST API, gRPC и Docker.

## 🚀 Функциональность

- Регистрация и вход (JWT + Cookie)
- Отправка выражений на вычисление
- Получение результатов и истории выражений
- gRPC-агенты для распределённых вычислений
- SQLite и GORM для хранения данных
- Интеграционные и модульные тесты

## 📦 Технологии

- Go, Gin, gRPC
- SQLite, GORM
- JWT, Cookie
- Docker, Docker Compose
- REST + gRPC API

## 🔧 Запуск

```bash
docker-compose up --build
````

* API доступен на `http://localhost:8080`
* gRPC-сервис: `:50051`

## 🔐 Авторизация

### Регистрация

```http
POST /api/v1/regin
Content-Type: application/json

{
  "login": "user1",
  "password": "pass123"
}
```

### Вход

```http
POST /api/v1/login
Content-Type: application/json

{
  "login": "user1",
  "password": "pass123"
}
```

## 📡 API

### Отправка выражения

```http
POST /api/v1/user/calculate
Cookie: jwt_access=<...>
Content-Type: application/json

{
  "expression": "1 + 2 * 3"
}
```

### Получение всех выражений

```http
GET /api/v1/user/expressions
Cookie: jwt_access=<...>
```

### Получение выражения по ID

```http
GET /api/v1/user/expressions/{id}
Cookie: jwt_access=<...>
```

## 🧪 Тестирование

```bash
go test ./...
```

## 📁 Переменные окружения

```env
PORT=8080
GRPC_PORT=50051
DATABASE_URL=sqlite.db
JWT_SECRET=your_jwt_secret
```
