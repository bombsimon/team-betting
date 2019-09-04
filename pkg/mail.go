package pkg

// MailService represents a service that can work with e-mail.
type MailService interface {
	SendMail(content *MailContent) error
}

// MailContent represents the data that can be added to a mail.
type MailContent struct {
	From    string
	To      string
	Subject string
	Body    string
}
