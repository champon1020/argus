package private

import (
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

var (
	errBadMultipartForm = errors.New("private.register_image: Cannot get multipart form from request context")
	errFailedOpen       = errors.New("private.register_image: Failed to open file")
	errFailedRead       = errors.New("private.register_image: Failed to read file")
	errFailedWrite      = errors.New("private.register_image: Failed to write file content")
)

// APIRegisterImage is the private handler to register new image.
func APIRegisterImage(ctx *gin.Context, _ model.DatabaseIface) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return argus.NewError(errBadMultipartForm, err)
	}

	fileHeaders := form.File["images"]
	dirPath := filepath.Join(argus.Env.Get("resource"), "images")

	// Save multipart form contents.
	for _, fh := range fileHeaders {
		body, err := getImageFromFileHeader(fh)
		if err != nil {
			return err
		}

		path := filepath.Join(dirPath, fh.Filename)

		if err := saveImage(path, body); err != nil {
			return err
		}
	}

	ctx.AbortWithStatus(http.StatusOK)
	return nil
}

func getImageFromFileHeader(fileHeader *multipart.FileHeader) ([]byte, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return nil, argus.NewError(errFailedOpen, err)
	}

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, argus.NewError(errFailedRead, err)
	}

	return body, nil
}

func saveImage(path string, body []byte) error {
	image, err := os.Create(path)
	if err != nil {
		return argus.NewError(errFailedOpen, err).
			AppendValue("path", path)
	}

	defer image.Close()
	if _, err := image.Write(body); err != nil {
		return argus.NewError(errFailedWrite, err)
	}

	return nil
}
