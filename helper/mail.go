package helper

import (
	"bytes"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"log"
	"os"
)

type Mail struct {
	from    string
	fromName string
	to      string
	toName string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func NewMail(to string, toName string, from string, fromName string, subject string) *Mail {
	return &Mail{
		to:      to,
		toName: toName,
		from: from,
		fromName: fromName,
		subject: subject,
	}
}

func (r *Mail) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Mail) sendMail() bool {
	from := mail.NewEmail(r.fromName, r.from)
	to := mail.NewEmail(r.toName, r.to)
	message := mail.NewSingleEmail(from, r.subject, to, "", r.body)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return false
	} else {
		fmt.Println(response.StatusCode)
		return true
	}
}

func (r *Mail) Send(templateName string, items interface{}) {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}
	if ok := r.sendMail(); ok {
		log.Println("SENT SUCCESS!!")
	} else {
		fmt.Println("SENT FAIL!!")
	}
}
