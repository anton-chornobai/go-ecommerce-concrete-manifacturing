package application

import (
	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
	"github.com/anton-chornobai/beton.git/internal/utils"
)

type UserService struct {
	repo domain.Repository
}

func NewUserService(repo domain.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(req domain.AuthenticationUserRequest) (domain.RegisterResult, error) {
	user, err := domain.NewUserCreated(req.Number)

	if err != nil {
		return domain.RegisterResult{}, err
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return domain.RegisterResult{}, err
	}

	if err := s.repo.Create(*user); err != nil {
		return domain.RegisterResult{}, err
	}

	return domain.RegisterResult{
		User:  *user,
		Token: token,
	}, nil
}

func (s *UserService) GetByPhone(number string) (domain.User, error) {
	return s.repo.GetByPhone(number)
}
