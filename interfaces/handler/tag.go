package handler

import (
	"github.com/champon1020/argus/config"
	"github.com/champon1020/argus/usecase"
	"github.com/labstack/echo/v4"
)

// TagHandler is interface that has handler functions related to tag.
type TagHandler interface {
	PublicTags(c echo.Context) error
}

type tagHandler struct {
	config *config.Config
	tu     usecase.TagUseCase
}

// NewTagHandler creates tagHandler.
func NewTagHandler(tu usecase.TagUseCase, config *config.Config) TagHandler {
	return &tagHandler{config: config, tu: tu}
}

func (th *tagHandler) PublicTags(c echo.Context) error {
	return nil
}
