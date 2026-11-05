package image_storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Storage struct {
	client *s3.Client
	bucket string
	cloudFrontURL string
}

func NewS3Storage(client *s3.Client, bucket string, cloudFrontURL string) *S3Storage {
	return &S3Storage{
		client: client,
		bucket: bucket,
		cloudFrontURL: strings.TrimRight(cloudFrontURL, "/"),
	}
}

func (s *S3Storage) Upload(
    ctx context.Context,
    file io.Reader,
    name string,
    contentType string,
) (string, error) {
    key := "products/" + uuid.NewString()

    _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        Body:        file,
        ContentType: aws.String(contentType),
    })
    if err != nil {
        return "", fmt.Errorf("failed to upload file to S3: %w", err)
    }

    imageURL := s.cloudFrontURL + "/" + key

    return imageURL, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
    _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return fmt.Errorf("failed to delete S3 object %q: %w", key, err)
    }

    return nil
}

func (s *S3Storage) ListObjects() {

}
