package service

import (
	"context"
	"errors"

	"testing"

	"github.com/anton-chornobai/beton.git/internal/modules/contact/domain"
	"github.com/anton-chornobai/beton.git/internal/modules/contact/dto"
)

type MockRepo struct {
	SavedContact *domain.UserContact
	PostCalled   bool

	DeleteCalled bool
	resErr       error
}

func (r *MockRepo) Save(ctx context.Context, uC *domain.UserContact) error {
	r.SavedContact = uC
	r.PostCalled = true
	return nil
}
func (r *MockRepo) Delete(ctx context.Context, id string) error {
	r.SavedContact = nil
	r.DeleteCalled = true
	return nil
}

func TestPost_InvalidNumber(t *testing.T) {
	repo := &MockRepo{}
	service := NewUserContactService(repo)

	zeroNumber := ""
	shortNumber := "096"
	tooLongNumber := "096096096096096096096096"
	invalidNumber := "+ABCD))9312"

	cases := []struct {
		TestName string
		Name        string
		Email       string
		Message     string
		Number      *string
		ThrownError error
	}{{
		TestName:    "Nil number",
		Name:        "Andy",
		Email:       "andy@gmail.com",
		Message:     "i would like to see the rest of the products",
		Number:      &zeroNumber,
		ThrownError: domain.ErrInvalidNumber,
	},
		{
			TestName: "Too Short Number",
			Name:        "Andy",
			Email:       "andy@gmail.com",
			Message:     "i would like to see the rest of the products",
			Number:      &shortNumber,
			ThrownError: domain.ErrInvalidNumber,
		},
		{
			TestName: "Too Long Number",
			Name:        "Andy",
			Email:       "andy@gmail.com",
			Message:     "i would like to see the rest of the products",
			Number:      &tooLongNumber,
			ThrownError: domain.ErrInvalidNumber,
		},
		{
			TestName: "Invalid Number Symbol",
			Name:        "Andy",
			Email:       "andy@gmail.com",
			Message:     "i would like to see the rest of the products",
			Number:      &invalidNumber,
			ThrownError: domain.ErrInvalidNumberSymbol,
		},
	}

	for _, tt := range cases {
		t.Run("Running case: " + tt.TestName, func(t *testing.T) {
			repo.PostCalled = false
			err := service.Post(context.Background(), dto.UserContactPostRequest{
				Name:    tt.Name,
				Email:   tt.Email,
				Message: tt.Message,
				Number:  tt.Number,
			})

			if err == nil {
				t.Fatalf("expected error %v, got nil", tt.ThrownError)
			}

			if !errors.Is(err, tt.ThrownError) {
				t.Fatalf("expected error %v got %v", tt.ThrownError, err)
			}

			if repo.PostCalled {
				t.Fatal("expected Save NOT to be called")
			}
		})
	}
}
