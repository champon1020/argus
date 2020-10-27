package private

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/handler"
	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

var (
	errFailedReadDir = errors.New("private.find_image: Failed to read directory")
	errBadPage       = errors.New("private.find_image: The page number parameter is bad")
)

// APIFindImagesRes is the response type.
type APIFindImagesRes struct {
	Images []string `json:"images"`
	Next   bool     `json:"next"`
}

// APIFindImages is the private handler to get all images.
func APIFindImages(ctx *gin.Context, _ model.DatabaseIface) error {
	// Channel for query parameter p.
	pCh := make(chan int, 1)

	// Channel for query parameter num.
	numCh := make(chan int, 1)

	// Channel for error variable.
	errCh := make(chan error, 2)

	// Response of this call.
	res := new(APIFindImagesRes)

	go handler.ParsePage(ctx, pCh, errCh)

	go handler.ParseNum(ctx, numCh, errCh)

	p, ok1 := <-pCh
	num, ok2 := <-numCh
	if !ok1 || !ok2 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errCh
	}

	dirPath := filepath.Join(argus.Env.Get("resource"), "images")
	images, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return argus.NewError(errFailedReadDir, err).
			AppendValue("dirPath", dirPath)
	}

	offset := (p - 1) * num
	if offset != 0 && offset >= len(images) {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return argus.NewError(errBadPage, nil).
			AppendValue("offset", offset).
			AppendValue("imagelen", len(images))
	}

	var fileInfo []os.FileInfo
	if len(images)-offset <= num {
		res.Next = false
		fileInfo = images[offset:len(images)]
	} else {
		res.Next = true
		fileInfo = images[offset : offset+num]
	}

	for _, fi := range fileInfo {
		res.Images = append(res.Images, fi.Name())
	}

	ctx.JSON(http.StatusOK, res)
	return nil
}
