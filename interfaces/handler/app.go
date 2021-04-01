package handler

import (
	"github.com/champon1020/argus/config"
	"github.com/champon1020/argus/usecase"
	"github.com/labstack/echo/v4"
)

// AppHandler includes all handlers.
type AppHandler interface {
	ArticleHandler
	TagHandler
	ImageHandler
}

type appHandler struct {
	config *config.Config
	aH     ArticleHandler
	tH     TagHandler
	iH     ImageHandler
}

// NewAppHandler creates appHandler.
func NewAppHandler(aU usecase.ArticleUseCase, tU usecase.TagUseCase, iU usecase.ImageUseCase, config *config.Config) AppHandler {
	return &appHandler{
		aH: NewArticleHandler(aU, tU, config),
		tH: NewTagHandler(tU, config),
		iH: NewImageHandler(iU, config),
	}
}

func (h *appHandler) PublicArticleByID(c echo.Context) error {
	return h.aH.PublicArticleByID(c)
}

func (h *appHandler) PublicArticles(c echo.Context) error {
	return h.aH.PublicArticles(c)
}

func (h *appHandler) PublicArticlesByTitle(c echo.Context) error {
	return h.aH.PublicArticlesByTitle(c)
}

func (h *appHandler) PublicArticlesByTag(c echo.Context) error {
	return h.aH.PublicArticlesByTag(c)
}

func (h *appHandler) ArticleByID(c echo.Context) error {
	return h.aH.ArticleByID(c)
}

func (h *appHandler) Articles(c echo.Context) error {
	return h.aH.Articles(c)
}

func (h *appHandler) PostArticle(c echo.Context) error {
	return h.aH.PostArticle(c)
}

func (h *appHandler) UpdateArticle(c echo.Context) error {
	return h.aH.UpdateArticle(c)
}

func (h *appHandler) UpdateArticleStatus(c echo.Context) error {
	return h.aH.UpdateArticleStatus(c)
}

func (h *appHandler) DeleteArticle(c echo.Context) error {
	return h.aH.DeleteArticle(c)
}

func (h *appHandler) PublicTags(c echo.Context) error {
	return h.tH.PublicTags(c)
}

func (h *appHandler) Images(c echo.Context) error {
	return h.iH.Images(c)
}

func (h *appHandler) PostImage(c echo.Context) error {
	return h.iH.PostImage(c)
}

func (h *appHandler) DeleteImage(c echo.Context) error {
	return h.iH.DeleteImage(c)
}
