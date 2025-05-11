# 1) Builder stage на Alpine с CGO и sqlite3-dev
FROM golang:1.24-alpine AS builder

# Устанавливаем компилятор и заголовки sqlite3
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем строго с CGO_ENABLED=1 и указанием GOOS=linux
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# 2) Финальный образ — чистый Alpine с рантайм-зависимостями
FROM alpine:latest

WORKDIR /app

# Устанавливаем только рантайм-библиотеки sqlite3 и certs
RUN apk add --no-cache sqlite-libs ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8000

CMD ["./main"]
