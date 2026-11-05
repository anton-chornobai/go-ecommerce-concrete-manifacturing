package mail

import (
    "context"
    "fmt"
    "os"
    "gopkg.in/gomail.v2"
)
func SendEmailTo(ctx context.Context, to, code string) error {
    if err := ctx.Err(); err != nil {
        return fmt.Errorf("email sending cancelled: %w", err)
    }

    password := os.Getenv("EMAIL_PASSWORD")
    if password == "" {
        return fmt.Errorf("EMAIL_PASSWORD is not set")
    }

    m := gomail.NewMessage()
    m.SetHeader("From", "antonchornobajj@gmail.com")
    m.SetHeader("To", to)
    m.SetHeader("Subject", "Код підтвердження")
    m.SetBody("text/html", "Введіть код підтвердження: <b>"+code+"</b>")

    d := gomail.NewDialer(
        "smtp.gmail.com",
        587,
        "antonchornobajj@gmail.com",
        password,
    )

    if err := d.DialAndSend(m); err != nil {
        return fmt.Errorf("send verification email: %w", err)
    }

    return nil
}