package usecase

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/champon1020/argus/domain/gcp"
	"github.com/champon1020/argus/usecase/pagenation"
	"github.com/champon1020/argus/util"
)

// ImageUseCase is usecase interface for image.
type ImageUseCase interface {
	ImageList(bktName string, p *pagenation.Pagenation) ([]string, error)
	CreateImage(file io.Reader, bktName, fileName string) error
	DeleteImage(bktName, imageURL string) error
}

type imageUseCase struct {
	cloudStorage gcp.CloudStorage
}

// NewImageUseCase creates imageUseCase.
func NewImageUseCase(cloudStorage gcp.CloudStorage) ImageUseCase {
	return &imageUseCase{cloudStorage: cloudStorage}
}

func (iU *imageUseCase) ImageList(bktName string, p *pagenation.Pagenation) ([]string, error) {
	ctx := context.Background()
	images, err := iU.cloudStorage.List(ctx, bktName)
	if err != nil {
		return nil, err
	}

	start := (p.Page - 1) * p.Limit
	end := util.MinInt(p.Page*p.Limit, len(images))
	p.Total = len(images)

	return images[start:end], nil
}

func (iU *imageUseCase) CreateImage(file io.Reader, bktName, fileName string) error {
	var img image.Image

	if filepath.Ext(fileName) == ".jpeg" || filepath.Ext(fileName) == ".jpg" {
		data, err := jpeg.Decode(file)
		if err != nil {
			return err
		}
		img = data
	} else if filepath.Ext(fileName) == ".png" {
		data, err := png.Decode(file)
		if err != nil {
			return err
		}
		img = data
	} else {
		return errors.New("invalid image format")
	}

	// Convert to webp.
	var dst bytes.Buffer
	if err := webp.Encode(&dst, img, nil); err != nil {
		return err
	}

	// Update the extention.
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".webp"

	ctx := context.Background()
	if err := iU.cloudStorage.Create(ctx, dst.Bytes(), bktName, fileName); err != nil {
		return err
	}

	return nil
}

func (iU *imageUseCase) DeleteImage(bktName, fileURL string) error {
	el := strings.Split(fileURL, bktName+"/images/")
	if len(el) == 1 {
		return errors.New("invalid image path")
	}

	ctx := context.Background()
	if err := iU.cloudStorage.Delete(ctx, bktName, el[1]); err != nil {
		return err
	}

	return nil
}
