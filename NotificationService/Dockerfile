FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY NotificationService/go.mod NotificationService/go.sum ./

RUN go mod download

COPY NotificationService/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o notificationservice ./cmd/main.go

FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/notificationservice .
COPY NotificationService/config/local.yaml ./config/local.yaml

CMD ["./notificationservice"]