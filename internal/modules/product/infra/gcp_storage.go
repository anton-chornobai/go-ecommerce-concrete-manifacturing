package infra

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

type GCPUploader struct {
	client *storage.Client
}

func NewGCPUploader(client *storage.Client) *GCPUploader {
	return &GCPUploader{
		client: client,
	}
}

func (s *GCPUploader) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	bucketName := "stezka_bucket"

	bucket := s.client.Bucket(bucketName)
	obj := bucket.Object(filename)
	writer := obj.NewWriter(ctx)

	if _, err := io.Copy(writer, file); err != nil {
		writer.Close()
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename), nil
}
