package application

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func (p *ProductService) cleanupImageKeys(keys []string) error {
    cleanupCtx, cancel := context.WithTimeout(
        context.Background(),
        10*time.Second,
    )
    defer cancel()

    var cleanupErr error

    for _, key := range keys {
        if err := p.imageStorage.Delete(cleanupCtx, key); err != nil {
            cleanupErr = errors.Join(
                cleanupErr,
                fmt.Errorf("delete S3 object %q: %w", key, err),
            )
        }
    }

    return cleanupErr
}