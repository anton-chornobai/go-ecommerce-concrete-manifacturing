package mail

import (
	"fmt"
	"net/smtp"
)


func SendEmailSample() {
	auth := smtp.PlainAuth(
		"",
		"antonchornobajj@gmail.com",
		"zsiq rlic adzi ftxu",
		"smtp.gmail.com",
	)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"antonchornobajj@gmail.com",
		[]string{"antonchornobajj@gmail.com"},
		[]byte("hello from golang"),
	)	
	if err != nil {
		fmt.Println("email sending error")
	}
}