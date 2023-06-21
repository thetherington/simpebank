package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error
}

type MailHogSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
	host              string
	port              string
}

func NewMailHogSender(name string, fromEmailAddress string, fromEmailPassword string, host string, port string) EmailSender {
	return &MailHogSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
		host:              host,
		port:              port,
	}
}

func (sender *MailHogSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()

	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	// smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, sender.host)

	// needed for mailhog when not localhost
	smtpAuth := smtp.CRAMMD5Auth(sender.fromEmailAddress, sender.fromEmailPassword)

	return e.Send(fmt.Sprintf("%s:%s", sender.host, sender.port), smtpAuth)
}
