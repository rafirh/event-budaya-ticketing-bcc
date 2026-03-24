package email

import (
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type MailjetSender struct {
	client   *mailjet.Client
	fromAddr string
	fromName string
}

type MailjetConfig struct {
	APIKey    string
	APISecret string
	FromAddr  string
	FromName  string
}

func NewMailjetSender(cfg MailjetConfig) *MailjetSender {
	client := mailjet.NewMailjetClient(cfg.APIKey, cfg.APISecret)

	return &MailjetSender{
		client:   client,
		fromAddr: cfg.FromAddr,
		fromName: cfg.FromName,
	}
}

func (s *MailjetSender) Send(to, subject, htmlBody string) error {
	if s.client == nil || s.fromAddr == "" {
		return fmt.Errorf("Mailjet config is incomplete")
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: s.fromAddr,
				Name:  s.fromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
				},
			},
			Subject:  subject,
			HTMLPart: htmlBody,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := s.client.SendMailV31(&messages)
	if err != nil {
		return fmt.Errorf("failed to send email with Mailjet: %w", err)
	}

	if res.ResultsV31[0].Status != "success" {
		return fmt.Errorf("Mailjet returned non-success status: %s", res.ResultsV31[0].Status)
	}

	return nil
}
