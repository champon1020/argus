package gcp

import (
	"context"
)

// CloudStorage is domain interface of cloud storage.
type CloudStorage interface {
	List(ctx context.Context, bktName, directory string) ([]string, error)
	Create(ctx context.Context, content []byte, bktName, filePath string) error
	Delete(ctx context.Context, bktName, filePath string) error
}
