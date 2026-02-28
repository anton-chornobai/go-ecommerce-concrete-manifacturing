package mail

import (
	"gopkg.in/gomail.v2"
	"os"
)

func SendEmailTo(to string, code string) error {
	password := os.Getenv("EMAIL_PASSWORD")
	if password == "" {
		panic("No email password is set")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "antonchornobajj@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Код підтвердження")
	m.SetBody("text/html", "Введіть код підтвердження: " + "<b>" +code+" </b>")
	d := gomail.NewDialer("smtp.gmail.com", 587, "antonchornobajj@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}
