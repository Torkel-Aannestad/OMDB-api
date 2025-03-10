package mailer

import (
	"bytes"
	"html/template"
	"time"

	"github.com/Torkel-Aannestad/OMDB-api/assets"
	"github.com/go-mail/mail/v2"
)

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 10 * time.Second
	dialer.StartTLSPolicy = mail.MandatoryStartTLS

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(assets.EmbededFiles, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	tmpl.ExecuteTemplate(subject, "subject", data)

	plainBody := new(bytes.Buffer)
	tmpl.ExecuteTemplate(plainBody, "plainBody", data)

	htmlBody := new(bytes.Buffer)
	tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	//mail retries, err == nil good case
	for i := 1; i < 3; i++ {
		err = m.dialer.DialAndSend(msg)
		if err == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return err
}
