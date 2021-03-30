package handler

import (
	"github.com/champon1020/argus/usecase"
	"github.com/labstack/echo/v4"
)

type AppHandler interface {
	ArticleHandler
}

type appHandler struct {
	ah ArticleHandler
}

func NewAppHandler(au usecase.ArticleUseCase) AppHandler {
	return &appHandler{
		ah: NewArticleHandler(au),
	}
}

func (h *appHandler) ArticleByID(c echo.Context) error {
	return h.ah.ArticleByID(c)
}
