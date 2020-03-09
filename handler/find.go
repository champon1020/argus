package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/champon1020/argus/service"

	"github.com/gin-gonic/gin"

	"github.com/champon1020/argus/repository"
)

type ResponseType struct {
	Articles []repository.Article `json:"articles"`
}

func FindArticleHandler(c *gin.Context) {
	var (
		articles []repository.Article
		response string
		argFlg   uint32
		err      error
	)

	mysql := repository.GlobalMysql
	argFlg = service.GenFlg(repository.Article{}, "Limit")
	if articles, err = repository.FindArticleCmd(mysql, repository.Article{}, argFlg); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
}

func FindArticleHandlerByTitle(c *gin.Context) {
	var (
		argArticle repository.Article
		articles   []repository.Article
		response   string
		argFlg     uint32
		err        error
	)

	argArticle.Title = c.Query("title")

	mysql := repository.GlobalMysql
	argFlg = service.GenFlg(repository.Article{}, "Title", "Limit")
	if articles, err = repository.FindArticleCmd(mysql, argArticle, argFlg); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
}

func FindArticleHandlerByCreateDate(c *gin.Context) {
	var (
		argArticle repository.Article
		articles   []repository.Article
		response   string
		argFlg     uint32
		err        error
	)

	if argArticle.CreateDate, err = time.Parse(time.RFC3339, c.Query("createDate")); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	argFlg = service.GenFlg(repository.Article{}, "CreateDate", "Limit")
	mysql := repository.GlobalMysql
	if articles, err = repository.FindArticleCmd(mysql, argArticle, argFlg); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
}

func FindArticleHandlerByCategory(c *gin.Context) {
	var (
		articles []repository.Article
		response string
		argFlg   uint32
		err      error
	)

	categoryNames := c.QueryArray("category")
	mysql := repository.GlobalMysql
	argFlg = service.GenFlg(repository.Article{}, "Limit")
	if articles, err = repository.FindArticleByCategoryCmd(mysql, categoryNames, argFlg); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
}

// Response type of Category
type CategoryResponseType struct {
	Categories []repository.CategoryResponse `json:"categories"`
}

func FindCategoryHandler(c *gin.Context) {
	var (
		categories []repository.CategoryResponse
		response   string
		argFlg     uint32
		err        error
	)

	mysql := repository.GlobalMysql
	argFlg = service.GenFlg(repository.Category{}, "Limit")
	if categories, err = repository.FindCategoryCmd(mysql, repository.Category{}, argFlg); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := CategoryResponseType{Categories: categories}
	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
}

// Response type of Draft
type DraftResponseType struct {
	Drafts []repository.Draft `json:"drafts"`
}

func FindDraftHandler(c *gin.Context) {
	var (
		drafts   []repository.Draft
		response string
		argFlg   uint32
		err      error
	)

	mysql := repository.GlobalMysql
	argFlg = service.GenFlg(repository.Draft{}, "Limit")
	if drafts, err = repository.FindDraftCmd(mysql, repository.Draft{}, argFlg); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := DraftResponseType{Drafts: drafts}
	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
}
