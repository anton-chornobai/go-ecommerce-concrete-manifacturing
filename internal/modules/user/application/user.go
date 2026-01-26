package application

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
	"github.com/anton-chornobai/beton.git/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo domain.Repository
	Logger  *slog.Logger
}

type EmailLoginRequest struct {
	Email    string
	Password string
}

func NewUserService(repo domain.Repository, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, Logger: logger}
}

func (s *UserService) Register(req domain.AuthenticationUserRequest) (domain.RegisterResult, error) {
	s.Logger.Debug("Register called", "number", req.Number)

	user, err := domain.CreateUser(req.Number)
	if err != nil {
		s.Logger.Error("Failed to create user domain object", "number", req.Number, "err", err)
		return domain.RegisterResult{}, err
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		s.Logger.Error("Failed to generate token", "err", err)
		return domain.RegisterResult{}, err
	}

	if err := s.repo.Create(*user); err != nil {
		s.Logger.Error("Failed to save user in repo", "number", user.Number, "err", err)
		return domain.RegisterResult{}, err
	}
	s.Logger.Info("User successfully registered", "number", user.Number)

	return domain.RegisterResult{
		User:  *user,
		Token: token,
	}, nil
}

func (s *UserService) SignUpByEmail(req EmailLoginRequest) (*domain.UserCreatedWithEmail, error) {
	const op = "UserService.SignUpByEmail"
	logger := s.Logger.With("op", op)

	if err := utils.ValidateEmailAndPassword(req.Email, req.Password); err != nil {
		logger.Warn("Invalid email credentials", "err", err)
		return nil, fmt.Errorf("signup failed: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", "err", err)
		return nil, fmt.Errorf("signup failed: %w", err)
	}


	user := domain.CreateUserWithEmail(req.Email, string(hashedPassword))

	if err := s.repo.SignUpByEmail(user); err != nil {
		logger.Error("Failed to save user in repo", "email", user.Email, "err", err)
		return nil, fmt.Errorf("signup failed: %w", err)
	}

	logger.Info("User successfully signed up", "user_id", user.ID, "email", user.Email)
	return user, nil
}

func (s *UserService) LoginByEmail(req EmailLoginRequest) (string, error) {
	const op = "UserService.LoginByEmail"

	if req.Email == "" || req.Password == "" {
		s.Logger.Warn("Empty credentials provided", "op", op)
		return "", errors.New("login failed: email and password required")
	}

	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		s.Logger.Warn("User not found", "op", op, "email", req.Email, "err", err)
		return "", fmt.Errorf("login failed: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.Logger.Warn("Invalid password attempt", "op", op, "email", req.Email)
		return "", errors.New("login failed: invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		s.Logger.Error("Failed to generate token", "op", op, "err", err)
		return "", fmt.Errorf("login failed: %w", err)
	}

	s.Logger.Info("Login successful", "op", op, "user_id", user.ID)
	return token, nil
}



func (s *UserService) GetByPhone(number string) (domain.User, error) {
	s.Logger.Debug("GetByPhone called", "number", number)
	user, err := s.repo.GetByPhone(number)
	if err != nil {
		s.Logger.Error("Failed to get user by phone", "number", number, "err", err)
		return domain.User{}, err
	}
	s.Logger.Debug("User retrieved by phone", "number", user.Number)
	return user, nil
}
