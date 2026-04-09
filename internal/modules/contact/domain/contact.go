package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const maxNameLength = 255
const maxEmailLength = 254

var (
	ErrNameTooLong   = fmt.Errorf("im'я занадто довге, більше %d знаків", maxNameLength)
	ErrEmailTooLong  = fmt.Errorf("емейл занадто довгий, максимально %d знаків", maxEmailLength)
	ErrInvalidNumber = fmt.Errorf("невірний номер телефону, має бути від 7 до 12 цифр")
	ErrInvalidNumerSymbol = fmt.Errorf("номер телефону містить недопустимі символи")
)

type UserContact struct {
	ID        string
	Name      string
	Email     string
	Number    string
	Message   string
	CreatedAt time.Time
}

func NewContact(name, email, number, message string) (*UserContact, error) {
	if len(name) > maxNameLength {
		return nil, ErrNameTooLong
	}
	if len(email) > maxNameLength {
		return nil, ErrEmailTooLong
	}
	if len(number) < 7 && len(number) > 12 {
		return nil, ErrInvalidNumber
	}
	for _, r := range number {
		if r < '0' || r > '9' {
			return nil, ErrInvalidNumerSymbol
		}
	}

	return &UserContact{
		ID:        uuid.NewString(),
		Name:      name,
		Email:     email,
		Number:    number,
		Message:   message,
		CreatedAt: time.Now(),
	}, nil
}

func (u *UserContact) Create() {

}

func (u *UserContact) Delete() {

}
