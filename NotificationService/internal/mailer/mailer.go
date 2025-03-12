package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
	Authorization   = "Authorization"
	Bearer          = "Bearer"
)

type EmailRequest struct {
	ToEmail string `json:"to_email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}



type EmailData struct {
	FromEmail string `json:"from_email"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Text      string `json:"text"`
}
type Mailer struct {
	ApiURL    string
	ApiToken  string
	FromEmail string
	http.Client
}

func NewMailer(ApiURL, ApiToken, FromEmail string) *Mailer {
	return &Mailer{ApiURL: ApiURL,
		ApiToken:  ApiToken,
		FromEmail: FromEmail,
	}
}

func (m *Mailer) SendEmail(ctx context.Context,toEmail, subject, body string) error {
	emailData := EmailData{
		FromEmail: m.FromEmail,
		To:        toEmail,
		Subject:   subject,
		Text:      body,
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		return fmt.Errorf("failed to marshal email data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, m.ApiURL+"/email/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create email request: %w", err)
	}

	req.Header.Set(ContentType, ApplicationJson)
	req.Header.Set(Authorization, Bearer+m.ApiToken)

	client := m.Client
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("failed to send email, status code: %d", resp.StatusCode)
	}

	return nil
}
