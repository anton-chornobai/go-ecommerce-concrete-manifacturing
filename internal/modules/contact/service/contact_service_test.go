package service

import (
	"context"
	"errors"
	"strings"

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
		TestName    string
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
			TestName:    "Too Short Number",
			Name:        "Andy",
			Email:       "andy@gmail.com",
			Message:     "i would like to see the rest of the products",
			Number:      &shortNumber,
			ThrownError: domain.ErrInvalidNumber,
		},
		{
			TestName:    "Too Long Number",
			Name:        "Andy",
			Email:       "andy@gmail.com",
			Message:     "i would like to see the rest of the products",
			Number:      &tooLongNumber,
			ThrownError: domain.ErrInvalidNumber,
		},
		{
			TestName:    "Invalid Number Symbol",
			Name:        "Andy",
			Email:       "andy@gmail.com",
			Message:     "i would like to see the rest of the products",
			Number:      &invalidNumber,
			ThrownError: domain.ErrInvalidNumberSymbol,
		},
	}

	for _, tt := range cases {
		t.Run("Running case:"+tt.TestName, func(t *testing.T) {
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

func TestPost_InvalidName(t *testing.T) {
	repo := &MockRepo{}
	service := NewUserContactService(repo)
	cases := []struct {
		TestName    string
		Name        string
		shouldThrow error
		Email       string
		Message     string
	}{
		{
			TestName:    "Name Too Long",
			Name:        strings.Repeat("a", 255),
			shouldThrow: domain.ErrNameTooLong,
			Email:       "test@mail.com",
			Message:     "hello",
		},
	}
	for _, tt := range cases {

		t.Run(tt.TestName, func(t *testing.T) {
			repo.PostCalled = false
			err := service.Post(context.Background(), dto.UserContactPostRequest{
				Name: tt.Name,
			})

			if err == nil {
				t.Fatalf("error is nil but should be %v", tt.shouldThrow)
			}

			if !errors.Is(err, tt.shouldThrow) {
				t.Fatalf("should throw: %v but got %v", tt.shouldThrow, err)
			}
			if repo.PostCalled {
				t.Fatal("expected Save NOT to be called")
			}
		})
	}
}

func TestPost_InvalidEmail(t *testing.T) {
	repo := &MockRepo{}
	service := NewUserContactService(repo)

	validNumber := "0971234567"

	cases := []struct {
		TestName    string
		Name        string
		Email       string
		Message     string
		Number      *string
		ThrownError error
	}{
		{
			TestName:    "Valid Email",
			Name:        "Andy",
			Email:       "test@mail.com",
			Message:     "hello",
			Number:      &validNumber,
			ThrownError: nil,
		},
		{
			TestName:    "Missing @",
			Name:        "Andy",
			Email:       "testmail.com",
			Message:     "hello",
			Number:      &validNumber,
			ThrownError: domain.ErrWrongEmailFormat,
		},
		{
			TestName:    "Too Long Email",
			Name:        "Andy",
			Email:       strings.Repeat("a", 256) + "@mail.com",
			Message:     "hello",
			Number:      &validNumber,
			ThrownError: domain.ErrEmailTooLong,
		},
	}

	for _, tt := range cases {
		t.Run("Running case:"+tt.TestName, func(t *testing.T) {
			repo.PostCalled = false

			err := service.Post(context.Background(), dto.UserContactPostRequest{
				Name:    tt.Name,
				Email:   tt.Email,
				Message: tt.Message,
				Number:  tt.Number,
			})

			if tt.ThrownError == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if !repo.PostCalled {
					t.Fatal("expected Save to be called")
				}

				return
			}

			if err == nil {
				t.Fatalf("expected error %v, got nil", tt.ThrownError)
			}

			if !errors.Is(err, tt.ThrownError) {
				t.Fatalf("expected %v, got %v", tt.ThrownError, err)
			}

			if repo.PostCalled {
				t.Fatal("expected Save NOT to be called")
			}
		})
	}
}
