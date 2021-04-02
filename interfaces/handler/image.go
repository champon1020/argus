package handler

import (
	"net/http"

	"github.com/champon1020/argus/config"
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/interfaces/handler/httputil"
	"github.com/champon1020/argus/usecase"
	"github.com/champon1020/argus/usecase/pagenation"
	"github.com/labstack/echo/v4"
)

// ImageHandler is handler ingerface for image.
type ImageHandler interface {
	Images(c echo.Context) error
	PostImage(c echo.Context) error
	DeleteImage(c echo.Context) error
}

type imageHandler struct {
	config *config.Config
	iU     usecase.ImageUseCase
}

// NewImageHandler creates imageHandler.
func NewImageHandler(iU usecase.ImageUseCase, config *config.Config) ImageHandler {
	return &imageHandler{iU: iU, config: config}
}

func (iH *imageHandler) Images(c echo.Context) error {
	page, err := httputil.ParsePage(c)
	if err != nil {
		// 400
		return err
	}

	p := &pagenation.Pagenation{Page: page, Limit: iH.config.LimitOnNumImages}

	images, err := iH.iU.ImageList(iH.config.StorageBucketName, p)
	if err != nil {
		// 500
		return err
	}

	return c.JSON(http.StatusOK, struct {
		ImageURLs  []string          `json:"image_urls"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{images, *p.MapToDomain()})
}

func (iH *imageHandler) PostImage(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		// 400
		return err
	}

	files := form.File["images"]
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		if err := iH.iU.CreateImage(src, iH.config.StorageBucketName, file.Filename); err != nil {
			return err
		}
	}

	return c.String(http.StatusOK, "success")
}

func (iH *imageHandler) DeleteImage(c echo.Context) error {
	url := c.QueryParam("url")

	if err := iH.iU.DeleteImage(iH.config.StorageBucketName, url); err != nil {
		return err
	}

	return c.String(http.StatusOK, "success")
}
