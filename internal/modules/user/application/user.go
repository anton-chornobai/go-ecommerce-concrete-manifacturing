package application

import (
	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
)

type UserAppService struct {
	service *domain.Service
}

func NewUserService(service *domain.Service) *UserAppService {
	return &UserAppService{
		service: service,
	}
}

func (r *UserAppService) Register(user domain.AuthenticationUserRequest) (domain.RegisterResult, error) {
	registeredUser, err := r.service.Register(user)

	if err != nil {
		return domain.RegisterResult{}, err
	}

	return  registeredUser, nil
}
