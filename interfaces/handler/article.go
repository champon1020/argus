package handler

import (
	"net/http"

	"github.com/champon1020/argus/config"
	"github.com/champon1020/argus/usecase"
	"github.com/labstack/echo/v4"
)

// ArticleHandler is handler interface for article.
type ArticleHandler interface {
	PublicArticleByID(c echo.Context) error
	PublicArticles(c echo.Context) error
	CountPublicArticles(c echo.Context) error
	ArticleByID(c echo.Context) error
	Articles(c echo.Context) error
	CountArticles(c echo.Context) error
	PostArticle(c echo.Context) error
	UpdateArticle(c echo.Context) error
	DeleteArticle(c echo.Context) error
}

type articleHandler struct {
	config *config.Config
	au     usecase.ArticleUseCase
}

// NewArticleHandler creates articleHandler.
func NewArticleHandler(au usecase.ArticleUseCase, config *config.Config) ArticleHandler {
	return &articleHandler{config: config, au: au}
}

func (h *articleHandler) PublicArticleByID(c echo.Context) error {
	return nil
}

func (h *articleHandler) PublicArticles(c echo.Context) error {
	return nil
}

func (h *articleHandler) CountPublicArticles(c echo.Context) error {
	return nil
}

func (h *articleHandler) ArticleByID(c echo.Context) error {
	id := c.Param("id")

	article, err := h.au.FindByID(h.config.DB, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, article)
}

func (h *articleHandler) Articles(c echo.Context) error {
	return nil
}

func (h *articleHandler) CountArticles(c echo.Context) error {
	return nil
}

func (h *articleHandler) PostArticle(c echo.Context) error {
	return nil
}

func (h *articleHandler) UpdateArticle(c echo.Context) error {
	return nil
}

func (h *articleHandler) DeleteArticle(c echo.Context) error {
	return nil
}
