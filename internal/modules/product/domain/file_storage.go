package domain

import (
	"context"
	"io"
)

type GCSUploader interface {
	Upload(ctx context.Context, file io.Reader, filename string) (string, error)
}	