package application

import (
	"fmt"

	"github.com/anton-chornobai/beton.git/internal/lib/jwt"
	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   domain.Repository
}

type EmailLoginRequest struct {
	Email    string
	Password string
}

func NewUserService(repo domain.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Signup(email, number string) (string, error) {

	user, err := domain.CreateUser(number)
	if err != nil {
		return "", err
	}

	token, err := jwtmanager.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	if err := s.repo.Create(*user); err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) SignUpByEmail(req EmailLoginRequest) (*domain.UserCreatedWithEmail, error) {

	// if err := jwtmanager.ValidateEmailAndPassword(req.Email, req.Password); err != nil {
	// 	logger.Warn("Invalid email credentials", "err", err)
	// 	return nil, fmt.Errorf("signup failed: %w", err)
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("signup failed: %w", err)
	}

	user := domain.CreateUserWithEmail(req.Email, string(hashedPassword))

	if err := s.repo.SignUpByEmail(user); err != nil {
		return nil, fmt.Errorf("signup failed: %w", err)
	}

	return user, nil
}

func (s *UserService) LoginByEmail(req EmailLoginRequest) (string, error) {

	if req.Email == "" || req.Password == "" {
		return "", errors.New("login failed: email and password required")
	}

	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("login failed: invalid credentials")
	}

	token, err := jwtmanager.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	return token, nil
}

func (s *UserService) GetByPhone(number string) (domain.User, error) {
	user, err := s.repo.GetByPhone(number)
	if err != nil {

		return domain.User{}, err
	}
	return user, nil
}
