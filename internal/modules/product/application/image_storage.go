package application

import (
	"context"
	"io"
)

type ImageStorage interface {
	Upload(ctx context.Context, file io.Reader, filename string) (string, error)
	Delete(ctx context.Context, imgURL string) error
}
