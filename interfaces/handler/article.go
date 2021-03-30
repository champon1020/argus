package handler

import (
	"github.com/champon1020/argus/usecase"
	"github.com/labstack/echo/v4"
)

type ArticleHandler interface {
	ArticleByID(c echo.Context) error
}

type articleHandler struct {
	articleUseCase usecase.ArticleUseCase
}

func NewArticleHandler(au usecase.ArticleUseCase) ArticleHandler {
	return &articleHandler{au}
}

func (h *articleHandler) ArticleByID(c echo.Context) error {
	return nil
}
