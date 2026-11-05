package application

import (
	"context"
	"time"

	"github.com/anton-chornobai/beton.git/internal/modules/user/infra/mail"
)

func (s *UserService) sendVerificationEmail(email, code string) {
    go func() {
        ctx, cancel := context.WithTimeout(
            context.Background(),
            15*time.Second,
        )
        defer cancel()

        if err := mail.SendEmailTo(ctx, email, code); err != nil {
            s.log.Error(
                "failed to send verification email",
                "email", email,
                "error", err,
            )
        }
    }()
}