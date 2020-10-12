package private

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/handler"
	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

// APIFindImagesRes is the response type.
type APIFindImagesRes struct {
	Images []string `json:"images"`
	Next   bool     `json:"next"`
}

// APIFindImages is the private handler to get all images.
func APIFindImages(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter p.
	pc := make(chan int, 1)

	// Channel for query parameter num.
	numc := make(chan int, 1)

	// Channel for error variable.
	errc := make(chan error, 2)

	// Response of this call.
	res := new(APIFindImagesRes)

	go handler.ParsePage(ctx, pc, errc)

	go handler.ParseNum(ctx, numc, errc)

	p, ok1 := <-pc
	num, ok2 := <-numc
	if !ok1 || !ok2 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errc
	}

	dirPath := filepath.Join(argus.Env.Get("resource"), "images")
	images, err := ioutil.ReadDir(dirPath)
	if err != nil {
		// handle error
	}

	offset := (p - 1) * num

	if offset != 0 && offset >= len(images) {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		// handle error
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
