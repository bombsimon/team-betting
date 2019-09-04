package mail

import (
	"crypto/tls"

	"github.com/bombsimon/team-betting/pkg"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

// Service represents a mail service.
type Service struct {
	Send func(...*gomail.Message) error
}

// New creates a new default service with a default dialer.
func New() *Service {
	d := gomail.NewDialer("localhost", 1025, "", "")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &Service{
		Send: d.DialAndSend,
	}
}

// SendMail will send an e-mail.
func (s *Service) SendMail(content *pkg.MailContent) error {
	m := gomail.NewMessage()
	m.SetHeader("From", content.From)
	m.SetHeader("To", content.To)
	m.SetHeader("Subject", content.Subject)
	m.SetBody("text/html", content.Body)

	if err := s.Send(m); err != nil {
		return errors.Wrap(err, "could not send email")
	}

	return nil
}
