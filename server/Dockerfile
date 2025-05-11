FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache sqlite-libs ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8000

CMD ["./main"]
