package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/anton-chornobai/beton.git/internal/mail"
	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
	"github.com/anton-chornobai/beton.git/pkg/utils"
)

type TokenManager interface {
	GenerateToken(id, email, role string) (string, error)
}

type VerificationaCodeManager interface {
	GenerateCode() (string, error)
	HashVerificationCode(string) string
	CompareHashAndCode(storedHash, userCode string) bool
}

type PasswordHasher interface {
	HashPassword(string) (string, error)
	CompareHashAndPassword(string, string) error
	HashVerificationCode(string) string
}

type UserService struct {
	repo                domain.Repository
	tokenManager        TokenManager
	passwordHasher      PasswordHasher
	verificationManager VerificationaCodeManager
	log                 *slog.Logger
}

func NewUserService(repo domain.Repository, tokenManager TokenManager, passwordHasher PasswordHasher, log *slog.Logger, verificationCodeManager VerificationaCodeManager) *UserService {
	return &UserService{
		repo:                repo,
		tokenManager:        tokenManager,
		passwordHasher:      passwordHasher,
		log:                 log,
		verificationManager: verificationCodeManager,
	}
}

func (s *UserService) SignupByEmail(ctx context.Context, email, password string) error {

	err := utils.ValidatePasswordAndEmail(email, password)

	if err != nil {
		return err
	}

	hashedPassword, err := s.passwordHasher.HashPassword(password)
	if err != nil {
		return fmt.Errorf("signup failed: %w", err)
	}

	user := domain.CreateUserWithEmail(email, hashedPassword)

	verificationCode, err := s.verificationManager.GenerateCode()

	if err != nil {
		s.log.Warn("failed to generate verification code")
		return fmt.Errorf("failed to generate verification code: %w", err)
	}
	hashedVerificationCode := s.verificationManager.HashVerificationCode(verificationCode)
	expiresAt := time.Now().UTC().Add(15 * time.Minute)

	if err := s.repo.SignupByEmail(ctx, user, hashedVerificationCode, &expiresAt); err != nil {
		return fmt.Errorf("signup failed: %w", err)
	}
	// Sending email aync
	go func(email, code string) {
		if err := mail.SendEmailTo(email, code); err != nil {
			s.log.Error("could not send verification code", "email", email, "err", err)
		}
	}(user.Email, verificationCode)

	return nil
}

func (s *UserService) VerifyUser(ctx context.Context, email, submittedCode string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)

	if user.VerificationExpiresAt == nil || user.VerificationExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("verification code expired %w", err)
	}
	if err != nil {
		return "", err
	}

	if user.IsVerified {
		return "", fmt.Errorf("user already verified")
	}

	if user.VerificationExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("verification code expired")
	}

	if !s.verificationManager.CompareHashAndCode(user.VerificationHash, submittedCode) {
		return "", fmt.Errorf("invalid verification code")
	}

	if err := s.repo.MarkUserVerified(ctx, email); err != nil {
		return "", fmt.Errorf("failed to verify user: %w", err)
	}

	token, err := s.tokenManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return token, nil
}

func (s *UserService) Signup(email, number string) (string, error) {
	user, err := domain.CreateUser(number)
	if err != nil {
		return "", err
	}

	token, err := s.tokenManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", err
	}

	if err := s.repo.Signup(user); err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) LoginByEmail(ctx context.Context, email, password string) (string, error) {
	err := utils.ValidatePasswordAndEmail(email, password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.LoginByEmail(ctx, email, password)

	if err != nil {
		return "", err
	}

	if err := s.passwordHasher.CompareHashAndPassword(user.Password, password); err != nil {
		return "", errors.New("login failed: invalid credentials")
	}

	token, err := s.tokenManager.GenerateToken(user.ID, user.Email, user.Role)

	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	return token, nil
}

func (s *UserService) GetByPhone(number string) (*domain.User, error) {
	user, err := s.repo.GetByPhone(number)
	if err != nil {

		return nil, err
	}
	return user, nil
}
