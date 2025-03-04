package mailer

import (
	"NotificationService/internal/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type EmailRequest struct {
	ToEmail string `json:"to_email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
type Mailer struct {
	ApiURL    string
	ApiToken  string
	FromEmail string
	logger    *logger.Logger
	http.Client
}

func NewMailer(ApiURL, ApiToken, FromEmail string, logger *logger.Logger) *Mailer {
	return &Mailer{ApiURL: ApiURL,
		ApiToken:  ApiToken,
		FromEmail: FromEmail,
		logger:    logger}
}

func (m *Mailer) SendEmail(ctx context.Context,toEmail, subject, body string) error {
	const (
		ContentType     = "Content-Type"
		ApplicationJson = "application/json"
		Authorization   = "Authorization"
		Bearer          = "Bearer"
	)

	m.logger.Info("Sending email",
		zap.String("to", toEmail),
		zap.String("subject", subject),
	)

	emailData := map[string]string{
		"from_email": m.FromEmail,
		"to":         toEmail,
		"subject":    subject,
		"text":       body,
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		m.logger.Error("failed to marshal email data", zap.Error(err))
		return fmt.Errorf("failed to marshal email data: %w", err)
	}

	req, err := http.NewRequest("POST", m.ApiURL+"/email/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		m.logger.Error("failed to create email request", zap.Error(err))
		return fmt.Errorf("failed to create email request: %w", err)
	}

	req.Header.Set(ContentType, ApplicationJson)
	req.Header.Set(Authorization, Bearer+m.ApiToken)

	client := m.Client
	resp, err := client.Do(req)
	if err != nil {
		m.logger.Error("failed to send email", zap.Error(err))
		return fmt.Errorf("failed to send email: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(resp.Body) 
		m.logger.Error("failed to send email",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response_body", string(body)),
		)
		return fmt.Errorf("failed to send email, status code: %d", resp.StatusCode)
	}

	return nil
}
