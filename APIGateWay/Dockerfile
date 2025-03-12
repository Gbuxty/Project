FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY proto/ /app/proto/

COPY APIGateWay/ ./APIGateWay/

WORKDIR /app/APIGateWay

RUN CGO_ENABLED=0 GOOS=linux go build -o apigateway ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/APIGateWay/apigateway .
COPY --from=builder /app/APIGateWay/swagger.yaml ./swagger.yaml

COPY APIGateWay/config/local.yaml ./config/local.yaml


EXPOSE 8080

ENTRYPOINT ["./apigateway","-c=config/local.yaml"]