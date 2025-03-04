package kafka

import (
	"NotificationService/internal/logger"
	"NotificationService/internal/mailer"
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader *kafka.Reader
	logger *logger.Logger
	mailer *mailer.Mailer
}

func NewConsumer(brokers []string, topic string, groupID string, mailer *mailer.Mailer, log *logger.Logger) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		}),
		logger: log,
		mailer: mailer,
	}
}

func (c *Consumer) Start(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            c.logger.Info("Kafka Consumer stopped")
            return
        default:
            msg, err := c.reader.ReadMessage(ctx)
            if err != nil {
                if err == context.Canceled {
                    c.logger.Info("Kafka Consumer context canceled")
                    return
                }
                c.logger.Error("Failed to read message from Kafka", zap.Error(err))
                time.Sleep(5 * time.Second)
                continue
            }

            var emailRequest mailer.EmailRequest
            if err := json.Unmarshal(msg.Value, &emailRequest); err != nil {
                c.logger.Error("Failed to unmarshal message", zap.Error(err))
                continue
            }

				if err := c.mailer.SendEmail(ctx, emailRequest.ToEmail, emailRequest.Subject, emailRequest.Body); err != nil {
					c.logger.Error("Failed to send email", zap.Error(err))
				} else {
                c.logger.Info("Email sent successfully", zap.String("to", emailRequest.ToEmail))
            }
        }
    }
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}