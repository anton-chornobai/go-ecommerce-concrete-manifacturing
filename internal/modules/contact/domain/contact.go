package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const maxNameLength = 255
const maxEmailLength = 254

var (
	ErrNameTooLong         = fmt.Errorf("im'я занадто довге, більше %d знаків", maxNameLength)
	ErrEmailTooLong        = fmt.Errorf("емейл занадто довгий, максимально %d знаків", maxEmailLength)
	ErrInvalidNumber       = fmt.Errorf("невірний номер телефону, має бути від 7 до 12 цифр")
	ErrInvalidNumberSymbol = fmt.Errorf("номер телефону містить недопустимі символи")
	ErrWrongEmailFormat    = errors.New("невірний формат емейлу")
)

type UserContact struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Message   string
	Number    *string
	CreatedAt time.Time
}

func NewContact(name, email, message string, number *string) (*UserContact, error) {
	if len(name) > maxNameLength {
		return nil, ErrNameTooLong
	}
	if len(email) > maxEmailLength {
		return nil, ErrEmailTooLong
	}
	if !strings.Contains(email, "@") {
		return nil, ErrWrongEmailFormat
	}

	if number != nil {
		if len(*number) < 7 || len(*number) > 12 {
			return nil, ErrInvalidNumber
		}
		for i, r := range *number {
			if i == 0 && r == '+' {
				continue
			}
			if r < '0' || r > '9' {
				return nil, ErrInvalidNumberSymbol
			}
		}
	}

	return &UserContact{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Number:    number,
		Message:   message,
		CreatedAt: time.Now(),
	}, nil
}
