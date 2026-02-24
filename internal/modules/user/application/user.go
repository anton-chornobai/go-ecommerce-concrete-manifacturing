package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
	"github.com/anton-chornobai/beton.git/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type TokenManager interface {
	GenerateToken(id, role string) (string, error)
}

type UserService struct {
	repo domain.Repository
	tokenManager TokenManager
}

func NewUserService(repo domain.Repository, tokenManager TokenManager) *UserService {
	return &UserService{repo: repo, tokenManager: tokenManager}
}

func (s *UserService) SignupByEmail(ctx context.Context, email, password string) (string, error) {

	err := utils.ValidatePasswordAndEmail(email, password)

	if err != nil {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("signup failed: %w", err)
	}

	user := domain.CreateUserWithEmail(email, string(hashedPassword))

	if err := s.repo.SignupByEmail(ctx, user); err != nil {
		return "", fmt.Errorf("signup failed: %w", err)
	}

	token, err := s.tokenManager.GenerateToken(user.ID, user.Role)

	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *UserService) Signup(email, number string) (string, error) {
	user, err := domain.CreateUser(number)
	if err != nil {
		return "", err
	}

	token, err := s.tokenManager.GenerateToken(user.ID, user.Role)
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

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return "", errors.New("login failed: invalid credentials")
	}

	token, err := s.tokenManager.GenerateToken(user.ID, user.Role)

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
