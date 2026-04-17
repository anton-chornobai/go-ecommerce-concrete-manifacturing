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

func strPtr(s string) *string {
	return &s
}

func TestPost_InvalidNumber(t *testing.T) {
	repo := &MockRepo{}
	service := NewUserContactService(repo)

	cases := []struct {
		TestName    string
		Number      string
		ThrownError error
	}{
		{"Empty number", "", domain.ErrInvalidNumber},
		{"Too short", "096", domain.ErrInvalidNumber},
		{"Too long", "096096096096096096096096", domain.ErrInvalidNumber},
		{"Invalid symbols", "+ABCD))9312", domain.ErrInvalidNumberSymbol},
	}

	for _, tt := range cases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo.PostCalled = false

			err := service.Post(context.Background(), dto.UserContactPostRequest{
				Name:    "Andy",
				Email:   strPtr("andy@gmail.com"),
				Message: "hello",
				Number:  tt.Number,
			})

			if err == nil {
				t.Fatalf("expected error %v, got nil", tt.ThrownError)
			}

			if !errors.Is(err, tt.ThrownError) {
				t.Fatalf("expected %v got %v", tt.ThrownError, err)
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

	longName := strings.Repeat("a", 256)

	err := service.Post(context.Background(), dto.UserContactPostRequest{
		Name:    longName,
		Email:   strPtr("test@mail.com"),
		Message: "hello",
		Number:  "0971234567",
	})

	if err == nil {
		t.Fatalf("expected error %v, got nil", domain.ErrNameTooLong)
	}

	if !errors.Is(err, domain.ErrNameTooLong) {
		t.Fatalf("expected %v got %v", domain.ErrNameTooLong, err)
	}

	if repo.PostCalled {
		t.Fatal("expected Save NOT to be called")
	}
}

func TestPost_InvalidEmail(t *testing.T) {
	repo := &MockRepo{}
	service := NewUserContactService(repo)

	validNumber := "0971234567"

	cases := []struct {
		TestName    string
		Email       *string
		ThrownError error
	}{
		{"Nil email (allowed)", nil, nil},
		{"Valid email", strPtr("test@mail.com"), nil},
		{"Missing @", strPtr("testmail.com"), domain.ErrWrongEmailFormat},
		{"Too long", strPtr(strings.Repeat("a", 256) + "@mail.com"), domain.ErrEmailTooLong},
	}

	for _, tt := range cases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo.PostCalled = false

			err := service.Post(context.Background(), dto.UserContactPostRequest{
				Name:    "Andy",
				Email:   tt.Email,
				Message: "hello",
				Number:  validNumber,
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