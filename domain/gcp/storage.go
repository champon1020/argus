package gcp

import (
	"context"
)

// CloudStorage is domain interface of cloud storage.
type CloudStorage interface {
	List(ctx context.Context, bktName string) ([]string, error)
	Create(ctx context.Context, content []byte, bktName, fileName string) error
	Delete(ctx context.Context, bktName, fileName string) error
}
