package email

import "net/smtp"

//go:generate moq -out emailtest/sender.go -pkg emailtest . Sender

type Sender interface {
	Send(from string, to []string, msg []byte) error
}

type SMTPSender struct {
	Addr string
	Auth smtp.Auth
}

func (S SMTPSender) Send(from string, to []string, msg []byte) error {
	return smtp.SendMail(
		S.Addr,
		S.Auth,
		from,
		to,
		msg)
}
