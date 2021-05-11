package handler

import (
	"io"
	"net/http"

	"github.com/champon1020/argus"
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
	HeaderImages(c echo.Context) error
	PostImage(c echo.Context) error
	DeleteImage(c echo.Context) error
}

type imageHandler struct {
	config *config.Config
	logger *argus.Logger
	iU     usecase.ImageUseCase
}

// NewImageHandler creates imageHandler.
func NewImageHandler(iU usecase.ImageUseCase, config *config.Config, logger *argus.Logger) ImageHandler {
	return &imageHandler{iU: iU, config: config, logger: logger}
}

// Images gets the image urls.
func (iH *imageHandler) Images(c echo.Context) error {
	page, err := httputil.ParsePage(c)
	if err != nil {
		iH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	p := &pagenation.Pagenation{Page: page, Limit: iH.config.LimitOnNumImages}

	images, err := iH.iU.ImageList(iH.config.StorageBucketName, p)
	if err != nil {
		iH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedGCSExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		ImageURLs  []string          `json:"image_urls"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{images, *p.MapToDomain()})
}

// HeaderImages gets the header image urls.
func (iH *imageHandler) HeaderImages(c echo.Context) error {
	images, err := iH.iU.HeaderImageList(iH.config.StorageBucketName)
	if err != nil {
		iH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedGCSExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		ImageURLs []string `json:"image_urls"`
	}{images})
}

// PostImage posts a new image.
func (iH *imageHandler) PostImage(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		iH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	files := form.File["image"]
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			iH.logger.Error(c, http.StatusInternalServerError, err)
			return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedOpenImage.Error())
		}
		defer src.Close()

		if err := iH.iU.CreateImage(src, iH.config.StorageBucketName, file.Filename); err != nil {
			iH.logger.Error(c, http.StatusInternalServerError, err)
			return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedGCSExec.Error())
		}
	}

	return c.String(http.StatusOK, "success")
}

// DeleteImage deletes the image.
func (iH *imageHandler) DeleteImage(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		iH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	if err := iH.iU.DeleteImages(iH.config.StorageBucketName, body); err != nil {
		iH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedGCSExec.Error())
	}

	return c.String(http.StatusOK, "success")
}
