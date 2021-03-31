package handler

import (
	"github.com/champon1020/argus/config"
	"github.com/champon1020/argus/usecase"
	"github.com/labstack/echo/v4"
)

// AppHandler includes all handlers.
type AppHandler interface {
	ArticleHandler
}

type appHandler struct {
	config *config.Config
	ah     ArticleHandler
}

// NewAppHandler creates appHandler.
func NewAppHandler(au usecase.ArticleUseCase, config *config.Config) AppHandler {
	return &appHandler{
		ah: NewArticleHandler(au, config),
	}
}

func (h *appHandler) PublicArticleByID(c echo.Context) error {
	return h.ah.PublicArticleByID(c)
}

func (h *appHandler) PublicArticles(c echo.Context) error {
	return h.ah.PublicArticles(c)
}

func (h *appHandler) CountPublicArticles(c echo.Context) error {
	return h.ah.CountPublicArticles(c)
}

func (h *appHandler) ArticleByID(c echo.Context) error {
	return h.ah.ArticleByID(c)
}

func (h *appHandler) Articles(c echo.Context) error {
	return h.ah.Articles(c)
}

func (h *appHandler) CountArticles(c echo.Context) error {
	return h.ah.CountArticles(c)
}

func (h *appHandler) PostArticle(c echo.Context) error {
	return h.ah.PostArticle(c)
}

func (h *appHandler) UpdateArticle(c echo.Context) error {
	return h.ah.UpdateArticle(c)
}

func (h *appHandler) DeleteArticle(c echo.Context) error {
	return h.ah.DeleteArticle(c)
}
