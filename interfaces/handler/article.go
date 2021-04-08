package handler

import (
	"errors"
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

func (aH *articleHandler) PublicArticleByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		aH.logger.Error(c, http.StatusBadRequest, errors.New("parameter id is empty"))
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	article, err := aH.aU.FindPublicByID(aH.config.DB, id)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Article domain.Article `json:"article"`
	}{*article})
}

func (aH *articleHandler) PublicArticles(c echo.Context) error {
	page, err := httputil.ParsePage(c)
	if err != nil {
		aH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	total, err := aH.aU.CountPublic(aH.config.DB)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	p := pagenation.NewPagenation(page, total, aH.config.LimitOnNumArticles)

	articles, err := aH.aU.FindPublic(aH.config.DB, *p)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Articles   []domain.Article  `json:"articles"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{*articles, *p.MapToDomain()})
}

func (aH *articleHandler) PublicArticlesByTitle(c echo.Context) error {
	title := c.Param("title")
	page, err := httputil.ParsePage(c)
	if err != nil {
		aH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	total, err := aH.aU.CountPublicByTitle(aH.config.DB, title)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	p := pagenation.NewPagenation(page, total, aH.config.LimitOnNumArticles)

	articles, err := aH.aU.FindPublicByTitle(aH.config.DB, *p, title)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Articles   []domain.Article  `json:"articles"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{*articles, *p.MapToDomain()})
}

func (aH *articleHandler) PublicArticlesByTag(c echo.Context) error {
	tag := c.Param("tag")
	page, err := httputil.ParsePage(c)
	if err != nil {
		aH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	total, err := aH.aU.CountPublicByTag(aH.config.DB, tag)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	p := pagenation.NewPagenation(page, total, aH.config.LimitOnNumArticles)

	articles, err := aH.aU.FindPublicByTag(aH.config.DB, *p, tag)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Articles   []domain.Article  `json:"articles"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{*articles, *p.MapToDomain()})
}

func (aH *articleHandler) ArticleByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		aH.logger.Error(c, http.StatusBadRequest, errors.New("path parameter id is empty"))
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	article, err := aH.aU.FindByID(aH.config.DB, id)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Article domain.Article `json:"article"`
	}{*article})
}

func (aH *articleHandler) Articles(c echo.Context) error {
	page, err := httputil.ParsePage(c)
	if err != nil {
		aH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	total, err := aH.aU.Count(aH.config.DB)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	p := pagenation.NewPagenation(page, total, aH.config.LimitOnNumPrivArticles)

	articles, err := aH.aU.Find(aH.config.DB, *p)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Articles   []domain.Article  `json:"articles"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{*articles, *p.MapToDomain()})
}

func (aH *articleHandler) DraftArticles(c echo.Context) error {
	page, err := httputil.ParsePage(c)
	if err != nil {
		aH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidQueryParam.Error())
	}

	total, err := aH.aU.CountDraftArticles(aH.config.DB)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	p := pagenation.NewPagenation(page, total, aH.config.LimitOnNumPrivArticles)

	articles, err := aH.aU.FindDraftArticles(aH.config.DB, *p)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Articles   []domain.Article  `json:"articles"`
		Pagenation domain.Pagenation `json:"pagenation"`
	}{*articles, *p.MapToDomain()})
}

func (aH *articleHandler) PostArticle(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		aH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	id, err := aH.aU.Post(aH.config.DB, body)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		ID string `json:"id"`
	}{id})
}

func (aH *articleHandler) UpdateArticle(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		aH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	id, err := aH.aU.Update(aH.config.DB, body)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		ID string `json:"id"`
	}{id})
}

func (aH *articleHandler) UpdateArticleStatus(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		aH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	id, err := aH.aU.UpdateStatus(aH.config.DB, body)
	if err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.JSON(http.StatusOK, struct {
		ID string `json:"id"`
	}{id})
}

func (aH *articleHandler) DeleteArticle(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		aH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	if err := aH.aU.Delete(aH.config.DB, body); err != nil {
		aH.logger.Error(c, http.StatusInternalServerError, err)
		return echo.NewHTTPError(http.StatusInternalServerError, ErrFailedDBExec.Error())
	}

	return c.String(http.StatusOK, "success")
}
