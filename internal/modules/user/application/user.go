package application

import (
	"log/slog"

	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
	"github.com/anton-chornobai/beton.git/internal/utils"
)

type UserService struct {
	repo   domain.Repository
	Logger *slog.Logger
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
	s.Logger.Debug("User domain object created", "role", user.Role)

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		s.Logger.Error("Failed to generate token", "err", err)
		return domain.RegisterResult{}, err
	}
	s.Logger.Debug("Token generated", "userID", user.ID)

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
