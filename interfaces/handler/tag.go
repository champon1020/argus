package handler

import (
	"net/http"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/config"
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/usecase"
	"github.com/labstack/echo/v4"
)

// TagHandler is interface that has handler functions related to tag.
type TagHandler interface {
	PublicTags(c echo.Context) error
}

type tagHandler struct {
	config *config.Config
	logger *argus.Logger
	tU     usecase.TagUseCase
}

// NewTagHandler creates tagHandler.
func NewTagHandler(tU usecase.TagUseCase, config *config.Config, logger *argus.Logger) TagHandler {
	return &tagHandler{config: config, logger: logger, tU: tU}
}

func (tH *tagHandler) PublicTags(c echo.Context) error {
	tags, err := tH.tU.FindPublic(tH.config.DB)
	if err != nil {
		tH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Tags []domain.Tag `json:"tags"`
	}{*tags})
}
