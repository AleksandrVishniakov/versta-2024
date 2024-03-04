package emailservice

import (
	"fmt"
	"net/smtp"
)

type EmailService interface {
	Write(content *EmailDTO) error
}

type emailService struct {
	auth     smtp.Auth
	sender   string
	hostname string
}

type EmailConfigs struct {
	Host        string
	Port        string
	SenderEmail string
	Password    string
}

func NewEmailService(cfg *EmailConfigs) EmailService {
	auth := smtp.PlainAuth(
		"",
		cfg.SenderEmail,
		cfg.Password,
		cfg.Host,
	)

	return &emailService{
		auth:     auth,
		sender:   cfg.SenderEmail,
		hostname: cfg.Host + ":" + cfg.Port,
	}
}

func (e *emailService) Write(content *EmailDTO) error {
	err := smtp.SendMail(
		e.hostname,
		e.auth,
		e.sender,
		[]string{content.To},
		formatEmailMessage(content.Subject, content.Body),
	)

	return err
}

func formatEmailMessage(subject, body string) []byte {
	message := fmt.Sprintf("Subject: %s\r\n%s", subject, body)

	return []byte(message)
}
