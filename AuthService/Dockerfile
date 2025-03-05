FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY AuthService/go.mod AuthService/go.sum ./
RUN go mod download

COPY AuthService/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o authservice ./cmd/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/authservice .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY AuthService/config/local.yaml ./config/local.yaml
COPY AuthService/migrations/postgres ./migrations


CMD goose -dir ./migrations postgres "user=postgres password=postgres_1234 dbname=Authentication sslmode=disable host=postgres" up && ./authservice