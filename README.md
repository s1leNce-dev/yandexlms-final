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

### Используем Docker Compose

```bash
docker compose up --build
````

* API будет доступен по адресу: `http://localhost:8000/api/v1`
* gRPC сервис будет запущен на порту: `:50051`

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
yandexlms-final/
│
├── server/                      # Серверная часть на Go
│   ├── app/                     # Инициализация приложения
│   │   ├── app.go
│   │   └── db.go                # Подключение к базе данных (SQLite + GORM)
│   ├── grpcservice/            # gRPC-сервер и обработчики задач
│   │   ├── grpcservice.go
│   │   ├── integration_test.go
│   │   └── server_test.go
│   ├── handlers/               # HTTP-обработчики (Gin)
│   │   ├── auth/               # Авторизация (JWT + Cookie)
│   │   │   └── auth.go
│   │   └── expressions/        # Вычисления (REST API)
│   │       └── expressions.go
│   ├── middlewares/           # JWT-миддлвар
│   │   └── authmiddleware.go
│   ├── models/                # GORM-модели (User, Expression, Task)
│   │   └── models.go
│   ├── proto/                 # gRPC-протокол
│   │   └── expression/
│   │       └── expression.proto
│   ├── routes/                # Маршруты Gin
│   │   └── routes.go
│   ├── utils/jwt/             # JWT-утилиты (генерация, парсинг)
│   │   └── jwt.go
│   ├── .env                   # Конфигурация окружения
│   ├── Dockerfile             # Docker-образ сервера
│   ├── go.mod
│   ├── go.sum
│   └── main.go                # Точка входа
│
├── agent/                     # gRPC-агент для выполнения задач
│   ├── eval/                  # Логика обработки выражений
│   │   └── eval.go
│   ├── proto/                 # Протокол gRPC (общий с сервером)
│   │   └── expression/
│   │       └── expression.proto
│   ├── .env                   # Переменные окружения агента
│   ├── agent.go               # Основной gRPC-клиент
│   ├── Dockerfile             # Docker-образ агента
│   ├── go.mod
│   └── go.sum
│
├── client/                    # Клиентская часть (React + Vite)
│   ├── public/                # Публичные ресурсы
│   ├── src/                   # Исходный код
│   │   ├── api/               # API-запросы (Fetch)
│   │   │   └── client.js
│   │   ├── assets/            # Статические ресурсы
│   │   ├── components/        # Компоненты интерфейса
│   │   │   └── ExpressionList.jsx
│   │   ├── context/           # Контекст авторизации
│   │   │   └── AuthContext.jsx
│   │   ├── pages/             # Страницы: вход, регистрация, главная
│   │   │   ├── HomePage.jsx / .css
│   │   │   ├── LoginPage.jsx
│   │   │   └── RegisterPage.jsx
│   │   ├── App.jsx / App.css  # Основное приложение
│   │   └── main.jsx           # Точка входа
│   ├── Dockerfile             # Docker-образ клиента
│   ├── index.html
│   ├── vite.config.js         # Конфигурация Vite
│   ├── package.json
│   └── README.md
│
├── docker-compose.yml         # Компоновка всех сервисов
└── README.md                  # Общая документация

```

