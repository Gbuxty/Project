FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем зависимости и скачиваем их
COPY NotificationService/go.mod NotificationService/go.sum ./
RUN go mod download

# Копируем исходный код
COPY NotificationService/ ./

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o notificationservice ./cmd/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник и конфиг
COPY --from=builder /app/notificationservice .
COPY NotificationService/config/local.yaml ./config/local.yaml

# Запуск приложения
CMD ["./notificationservice"]