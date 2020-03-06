package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/champon1020/argus/repository"
)

type ResponseType struct {
	Articles []repository.Article `json:"articles"`
}

type CategoryResponseType struct {
	Categories []repository.CategoryResponse `json:"categories"`
}

func FindArticleHandler(c *gin.Context) {
	var (
		err      error
		articles []repository.Article
		response string
	)

	mysql := repository.GlobalMysql
	if articles, err = repository.FindArticleCmd(mysql, repository.Article{}, 0); err != nil {
		(c.Writer).WriteHeader(http.StatusInternalServerError)
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = ParseToJson(&res); err != nil {
		(c.Writer).WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
}

func FindArticleHandlerByTitle(c *gin.Context) {
	var (
		argArticle repository.Article
		articles   []repository.Article
		response   string
		err        error
		argFlg     uint32
	)

	argArticle.Title = c.Query("title")

	argFlg = GenFlg(repository.Article{}, "Title")
	mysql := repository.GlobalMysql
	if articles, err = repository.FindArticleCmd(mysql, argArticle, argFlg); err != nil {
		(c.Writer).WriteHeader(http.StatusInternalServerError)
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = ParseToJson(&res); err != nil {
		(c.Writer).WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
}

func FindArticleHandlerByCreateDate(c *gin.Context) {
	var (
		argArticle repository.Article
		articles   []repository.Article
		response   string
		err        error
		argFlg     uint32
	)

	if argArticle.CreateDate, err = time.Parse(time.RFC3339, c.Query("createDate")); err != nil {
		(c.Writer).WriteHeader(http.StatusInternalServerError)
		return
	}

	argFlg = GenFlg(repository.Article{}, "create_date")
	mysql := repository.GlobalMysql
	if articles, err = repository.FindArticleCmd(mysql, argArticle, argFlg); err != nil {
		(c.Writer).WriteHeader(http.StatusInternalServerError)
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = ParseToJson(&res); err != nil {
		(c.Writer).WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
}

func FindArticleHandlerByCategory(c *gin.Context) {
	var (
		articles []repository.Article
		response string
		err      error
	)

	categoryNames := c.QueryArray("category")
	mysql := repository.GlobalMysql
	if articles, err = repository.FindArticleByCategoryCmd(mysql, categoryNames); err != nil {
		(c.Writer).WriteHeader(http.StatusInternalServerError)
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = ParseToJson(&res); err != nil {
		(c.Writer).WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
}

func FindCategoryHandler(c *gin.Context) {
	var (
		err        error
		categories []repository.CategoryResponse
		response   string
	)

	mysql := repository.GlobalMysql
	if categories, err = repository.FindCategoryCmd(mysql, repository.Category{}, 0); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	res := CategoryResponseType{Categories: categories}
	if response, err = ParseToJson(&res); err != nil {
		(c.Writer).WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
}
