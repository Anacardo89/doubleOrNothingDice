package email

import (
	"fmt"
	"net/smtp"
	"strconv"
)

type EmailSender struct {
	SMTPHost       string
	SMTPPort       int
	SenderEmail    string
	SenderPassword string
}

func NewEmailSender(smtpHost string, smtpPort int, senderEmail, senderPassword string) *EmailSender {
	return &EmailSender{
		SMTPHost:       smtpHost,
		SMTPPort:       smtpPort,
		SenderEmail:    senderEmail,
		SenderPassword: senderPassword,
	}
}

func (es *EmailSender) Send(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", es.SenderEmail, es.SenderPassword, es.SMTPHost)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
	addr := es.SMTPHost + ":" + strconv.Itoa(es.SMTPPort)
	return smtp.SendMail(addr, auth, es.SenderEmail, []string{to}, msg)
}
