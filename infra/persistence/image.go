package persistence

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"github.com/champon1020/argus/domain/gcp"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type imagePersistence struct{}

// NewImagePersistence is persistence interface for image.
func NewImagePersistence() gcp.CloudStorage {
	return &imagePersistence{}
}

func (iP *imagePersistence) List(ctx context.Context, bktName string) ([]string, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("ARGUS_CLOUD_STORAGE_KEY_PATH")))
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0, 10)
	it := client.Bucket(bktName).Objects(ctx, &storage.Query{Prefix: "images/"})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		urls = append(urls, fmt.Sprintf("https://storage.googleapis.com/%s/%s", bktName, attrs.Name))
	}

	return urls[1:], nil
}

func (iP *imagePersistence) Create(ctx context.Context, content []byte, bktName, fileName string) error {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("ARGUS_CLOUD_STORAGE_KEY_PATH")))
	if err != nil {
		return err
	}

	wc := client.Bucket(bktName).Object("images/" + fileName).NewWriter(ctx)
	wc.ContentType = "image/webp"
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	if _, err := wc.Write(content); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (iP *imagePersistence) Delete(ctx context.Context, bktName, fileName string) error {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("ARGUS_CLOUD_STORAGE_KEY_PATH")))
	if err != nil {
		return err
	}

	wc := client.Bucket(bktName).Object("images/" + fileName)
	if err := wc.Delete(ctx); err != nil {
		return err
	}

	return nil
}
