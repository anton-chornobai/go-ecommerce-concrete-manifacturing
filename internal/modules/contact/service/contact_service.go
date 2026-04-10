package service

import (
	"context"
	"github.com/anton-chornobai/beton.git/internal/modules/contact/domain"
	"github.com/anton-chornobai/beton.git/internal/modules/contact/dto"
)

type UserContactService struct {
	repo domain.Repository
}

func NewUserContactService(repo domain.Repository) *UserContactService {
	return &UserContactService{
		repo: repo,
	}
}

func (u *UserContactService) Post(ctx context.Context, req dto.UserContactPostRequest) error {
	userContact, err := domain.NewContact(req.Name, req.Email, req.Message, req.Number)
	if err != nil {
		return err
	}

	err = u.repo.Save(ctx, userContact)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserContactService) Delete(ctx context.Context, id string) error {
	err := u.repo.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
