package main

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type Mailer interface {
	Send(to, subject, body string) error
}

type mailer struct {
	account  string
	password string
	server   string
	host     string
}

func New(server, account, password string) (Mailer, error) {
	host := strings.Split(server, ":")
	if len(host) != 2 {
		return nil, errors.New("server address format: <HOST>:<PORT>")
	}

	return &mailer{
		account:  account,
		password: password,
		server:   server,
		host:     host[0],
	}, nil
}

func (m *mailer) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.account, m.password, m.host)
	msg := "From: " + m.account + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	err := smtp.SendMail(m.server, auth, m.account, []string{to}, []byte(msg))
	if err != nil {
		fmt.Errorf("%w", err)
	}
	return nil

}

func main() {
	email := os.Getenv("EMAIL")
	mailAccount := strings.Split(email, "/")
	if len(mailAccount) != 3 {
		log.Println("Environmental vairable 'EMAIL' is required and format is <smtp server>:<port>/<email>/<password>")
		return
	}
	mailer, err := New(mailAccount[0], mailAccount[1], mailAccount[2])
	if err != nil {
		log.Fatal(err)
	}

	if err := mailer.Send(os.Args[1], os.Args[2], "HELLO!"); err != nil {
		log.Fatal(err)
	}
}
