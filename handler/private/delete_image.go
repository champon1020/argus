package private

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

var (
	errNoSuchImageFile = errors.New("private.delete_image.go: No such image file")
	errFailedRemove    = errors.New("private.delete_image.go: Failed to remove image file")
)

// APIDeleteImageReq is the request type.
type APIDeleteImageReq struct {
	ImageNames []string `json"imageNames"`
}

// APIDeleteImage is the private handler to delete image.
func APIDeleteImage(ctx *gin.Context, _ model.DatabaseIface) error {
	// Channel for request.
	reqc := make(chan APIDeleteImageReq, 1)

	// Channel for error variable.
	errc := make(chan error, 1)

	go ParseDeleteImage(ctx, reqc, errc)

	req, ok := <-reqc
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errc
	}

	flgDelete := true

	wg := new(sync.WaitGroup)
	for _, image := range req.ImageNames {
		wg.Add(1)
		go func() {
			defer wg.Done()
			path := filepath.Join(argus.Env.Get("resource"), "images", image)
			if err := deleteImage(path); err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				flgDelete = false
				errc <- err
				return
			}
		}()
	}

	wg.Wait()

	if !flgDelete {
		return <-errc
	}

	ctx.AbortWithStatus(http.StatusOK)
	return nil
}

func deleteImage(path string) error {
	if _, err := os.Stat(path); err != nil {
		return argus.NewError(errNoSuchImageFile, err).
			AppendValue("path", path)
	}

	if err := os.Remove(path); err != nil {
		return argus.NewError(errFailedRemove, err).
			AppendValue("path", path)
	}

	return nil
}
