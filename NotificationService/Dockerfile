FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY NotificationService/ ./NotificationService/

WORKDIR /app/NotificationService

RUN CGO_ENABLED=0 GOOS=linux go build -o notificationservice ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/NotificationService/notificationservice .

COPY NotificationService/config/local.yaml ./config/local.yaml

ENTRYPOINT ["./notificationservice", "-c=config/local.yaml"]