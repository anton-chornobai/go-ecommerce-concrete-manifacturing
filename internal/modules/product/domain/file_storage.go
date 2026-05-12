package domain

import (
	"context"
	"io"
	"mime/multipart"
)

type GCSUploader interface {
	Upload(ctx context.Context, file io.Reader, filename string) (string, error)
	UploadMultipleFiles(ctx context.Context, headers []*multipart.FileHeader) ([]string, error)
	Delete(ctx context.Context, imgURL string) error
}
