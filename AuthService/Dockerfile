FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY proto/ /app/proto/

COPY AuthService/ ./AuthService/

WORKDIR /app/AuthService

RUN CGO_ENABLED=0 GOOS=linux go build -o authservice ./cmd/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/AuthService/authservice .

COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY AuthService/config/local.yaml ./config/local.yaml
COPY AuthService/migrations/postgres ./migrations

EXPOSE 9090

ENTRYPOINT ["./authservice", "-c=config/local.yaml"]