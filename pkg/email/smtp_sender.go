package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type Sender interface {
	Send(to, subject, htmlBody string) error
}

type SMTPSender struct {
	host       string
	port       string
	username   string
	password   string
	encryption string
	fromAddr   string
	fromName   string
}

type SMTPConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Encryption string
	FromAddr   string
	FromName   string
}

func NewSMTPSender(cfg SMTPConfig) *SMTPSender {
	fromAddr := strings.TrimSpace(cfg.FromAddr)
	if fromAddr == "" {
		fromAddr = cfg.Username
	}

	return &SMTPSender{
		host:       cfg.Host,
		port:       cfg.Port,
		username:   cfg.Username,
		password:   cfg.Password,
		encryption: strings.ToLower(strings.TrimSpace(cfg.Encryption)),
		fromAddr:   fromAddr,
		fromName:   cfg.FromName,
	}
}

func (s *SMTPSender) Send(to, subject, htmlBody string) error {
	if s.host == "" || s.port == "" || s.username == "" || s.password == "" || s.fromAddr == "" {
		return fmt.Errorf("SMTP config is incomplete")
	}

	header := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", s.fromName, s.fromAddr),
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	var msg strings.Builder
	for k, v := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	if s.encryption == "ssl" {
		return s.sendWithTLS(addr, to, msg.String(), auth)
	}

	return smtp.SendMail(addr, auth, s.fromAddr, []string{to}, []byte(msg.String()))
}

func (s *SMTPSender) sendWithTLS(addr, to, message string, auth smtp.Auth) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: s.host})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(s.fromAddr); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := writer.Write([]byte(message)); err != nil {
		_ = writer.Close()
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return client.Quit()
}
