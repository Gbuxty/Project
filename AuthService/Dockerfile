FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем зависимости и скачиваем их
COPY AuthService/go.mod AuthService/go.sum ./
RUN go mod download

# Копируем исходный код
COPY AuthService/ ./

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o authservice ./cmd/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник и конфиг
COPY --from=builder /app/authservice .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY AuthService/config/local.yaml ./config/local.yaml
COPY AuthService/migrations/postgres ./migrations

# Запуск приложения
CMD goose -dir ./migrations postgres "user=postgres password=postgres_1234 dbname=Authentication sslmode=disable host=postgres" up && ./authservice