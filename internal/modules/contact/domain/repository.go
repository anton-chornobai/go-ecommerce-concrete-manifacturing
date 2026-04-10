package domain

import "context"

type Repository interface {
	Save(ctx context.Context, userContact *UserContact) error
	Delete(ctx context.Context, id string) error
}
