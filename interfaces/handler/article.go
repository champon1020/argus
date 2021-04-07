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

// ArticleHandler is handler interface for article.
type ArticleHandler interface {
	PublicArticleByID(c echo.Context) error
	PublicArticles(c echo.Context) error
	PublicArticlesByTitle(c echo.Context) error
	PublicArticlesByTag(c echo.Context) error
	ArticleByID(c echo.Context) error
	Articles(c echo.Context) error
	DraftArticles(c echo.Context) error
	PostArticle(c echo.Context) error
	UpdateArticle(c echo.Context) error
	UpdateArticleStatus(c echo.Context) error
	DeleteArticle(c echo.Context) error
}

type articleHandler struct {
	config *config.Config
	logger *argus.Logger
	aU     usecase.ArticleUseCase
	tU     usecase.TagUseCase
}

// NewArticleHandler creates articleHandler.
func NewArticleHandler(aU usecase.ArticleUseCase, tU usecase.TagUseCase, config *config.Config, logger *argus.Logger) ArticleHandler {
	return &articleHandler{config: config, logger: logger, aU: aU, tU: tU}
}

func (h *articleHandler) PublicArticleByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	article, err := h.aU.FindPublicByID(h.config.DB, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Article domain.Article `json:"article"`
	}{*article})
}

func (h *articleHandler) PublicArticles(c echo.Context) error {
	page, err := httputil.ParsePage(c)
	if err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	total, err := h.aU.CountPublic(h.config.DB)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	p := pagenation.NewPagenation(page, total, h.config.LimitOnNumArticles)

	articles, err := h.aU.FindPublic(h.config.DB, *p)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Articles   []domain.Article  `json:"articles"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{*articles, *p.MapToDomain()})
}

func (h *articleHandler) PublicArticlesByTitle(c echo.Context) error {
	title := c.Param("title")
	page, err := httputil.ParsePage(c)
	if err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	total, err := h.aU.CountPublicByTitle(h.config.DB, title)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	p := pagenation.NewPagenation(page, total, h.config.LimitOnNumArticles)

	articles, err := h.aU.FindPublicByTitle(h.config.DB, *p, title)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Articles   []domain.Article  `json:"articles"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{*articles, *p.MapToDomain()})
}

func (h *articleHandler) PublicArticlesByTag(c echo.Context) error {
	tag := c.Param("tag")
	page, err := httputil.ParsePage(c)
	if err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	total, err := h.aU.CountPublicByTag(h.config.DB, tag)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	p := pagenation.NewPagenation(page, total, h.config.LimitOnNumArticles)

	articles, err := h.aU.FindPublicByTag(h.config.DB, *p, tag)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Articles   []domain.Article  `json:"articles"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{*articles, *p.MapToDomain()})
}

func (h *articleHandler) ArticleByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	article, err := h.aU.FindByID(h.config.DB, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Article domain.Article `json:"article"`
	}{*article})
}

func (h *articleHandler) Articles(c echo.Context) error {
	page, err := httputil.ParsePage(c)
	if err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	total, err := h.aU.Count(h.config.DB)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	p := pagenation.NewPagenation(page, total, h.config.LimitOnNumPrivArticles)

	articles, err := h.aU.Find(h.config.DB, *p)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Articles   []domain.Article  `json:"articles"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{*articles, *p.MapToDomain()})
}

func (h *articleHandler) DraftArticles(c echo.Context) error {
	page, err := httputil.ParsePage(c)
	if err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	total, err := h.aU.CountDraftArticles(h.config.DB)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	p := pagenation.NewPagenation(page, total, h.config.LimitOnNumPrivArticles)

	articles, err := h.aU.FindDraftArticles(h.config.DB, *p)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Articles   []domain.Article  `json:"articles"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{*articles, *p.MapToDomain()})
}

func (h *articleHandler) PostArticle(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	id, err := h.aU.Post(h.config.DB, body)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		ID string `json:"id"`
	}{id})
}

func (h *articleHandler) UpdateArticle(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	id, err := h.aU.Update(h.config.DB, body)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		ID string `json:"id"`
	}{id})
}

func (h *articleHandler) UpdateArticleStatus(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	id, err := h.aU.UpdateStatus(h.config.DB, body)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		ID string `json:"id"`
	}{id})
}

func (h *articleHandler) DeleteArticle(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	if err := h.aU.Delete(h.config.DB, body); err != nil {
		return c.String(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.String(http.StatusOK, "success")
}
