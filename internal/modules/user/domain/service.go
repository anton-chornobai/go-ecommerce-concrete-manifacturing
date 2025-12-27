package domain

import (
	"errors"

	"github.com/anton-chornobai/beton.git/internal/utils"
	"github.com/google/uuid"
)

const phoneNumberLength = 9

type Service struct {
	repo Repository
}

type RegisterResult struct {
    User  User
    Token string
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(user AuthenticationUserRequest) (RegisterResult, error) {
	if user.Number == "" {
		return  RegisterResult{}, errors.New("phone number is required")
	}

	createdUser := User {
		ID:     uuid.NewString(),
		Role:   "user",
		Number: user.Number,
	}

	token, err := utils.GenerateToken(createdUser.ID, createdUser.Role)
	if err != nil {
		return RegisterResult{}, err
	}

	err = s.repo.Create(createdUser)

	if err != nil {
		return RegisterResult{}, err
	}

	return RegisterResult{
		User: createdUser,
		Token: token,
	}, nil
}

func (s *Service) GetByPhone(number string) (User, error) {
	if number == "" || len(number) < phoneNumberLength {
		return User{}, errors.New("invalid number argument")
	}
	
	user, err := s.repo.GetByPhone(number) 

	if err != nil {
		return User{}, err
	}
	
	return user, nil
}