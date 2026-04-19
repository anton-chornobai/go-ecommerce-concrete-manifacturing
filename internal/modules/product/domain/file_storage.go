package domain

import (
	"context"
	"io"
)

type FileStorage interface {
	Upload(ctx context.Context, file io.Reader, filename string) (string, error)
}	