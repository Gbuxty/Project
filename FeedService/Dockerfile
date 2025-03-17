FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY proto/ /app/proto/

COPY FeedService/ ./FeedService/

WORKDIR /app/FeedService

RUN CGO_ENABLED=0 GOOS=linux go build -o feedservice ./cmd/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/FeedService/feedservice .

COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY FeedService/config/local.yaml ./config/local.yaml
COPY FeedService/migrations/postgres ./migrations

EXPOSE 7070

ENTRYPOINT ["./feedservice", "-c=config/local.yaml"]