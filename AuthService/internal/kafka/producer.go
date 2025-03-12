package kafka

import (
	"Project/AuthService/internal/logger"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	logger *logger.Logger
}

func NewProducer(brokers []string, topic string, logger *logger.Logger) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		logger: logger,
	}
}

func (p *Producer) SendMessage(ctx context.Context, key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: jsonValue,
	})

	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}
	p.logger.Info("Message sent to Kafka", zap.String("key", key))
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
