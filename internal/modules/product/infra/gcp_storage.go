package infra

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

func (s *GCPUploader) UploadMultipleFiles(ctx context.Context, headers []*multipart.FileHeader) ([]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type result struct {
		index int
		url   string
		err   error
	}

	resultChan := make(chan result, len(headers))
	var wg sync.WaitGroup

	for i, header := range headers {
		wg.Add(1)
		go func(i int, h *multipart.FileHeader) {
			defer wg.Done()

			f, err := h.Open()
			if err != nil {
				resultChan <- result{index: i, err: err}
				return
			}
			defer f.Close()

			ext := filepath.Ext(h.Filename)
			nameOnly := strings.TrimSuffix(h.Filename, ext)
			uniqueName := fmt.Sprintf("%s_%d_%d%s", nameOnly, time.Now().UnixNano(), i, ext)

			url, err := s.Upload(ctx, f, uniqueName)
			if err != nil {
				resultChan <- result{index: i, err: err}
				return
			}

			resultChan <- result{index: i, url: url}
		}(i, header)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	urls := make([]string, len(headers))
	for res := range resultChan {
		if res.err != nil {
			cancel()
			return nil, fmt.Errorf("не вдалося завантажити файл: %w", res.err)
		}
		urls[res.index] = res.url
	}

	return urls, nil
}

func (s *GCPUploader) Delete(ctx context.Context, url string) error {
	bucketName := "stezka_bucket"

	prefix := fmt.Sprintf("https://storage.googleapis.com/%s/", bucketName)
	filename := strings.TrimPrefix(url, prefix)
	if filename == url {
		return fmt.Errorf("невалідний url зображення: %s", url)
	}

	obj := s.client.Bucket(bucketName).Object(filename)
	if err := obj.Delete(ctx); err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			// object already gone — not an error from our perspective
			return nil
		}
		return fmt.Errorf("не вдалося видалити файл %s: %w", filename, err)
	}

	return nil
}
