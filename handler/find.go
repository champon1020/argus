package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/champon1020/argus/repository"
)

type QueryParam struct {
	Article repository.Article
}

type ResponseType struct {
	Articles []repository.Article `json:"articles"`
}

type CategoryResponseType struct {
	Categories []repository.Category `json:"categories"`
}

func FindArticleHandler(c *gin.Context) {
	var (
		err      error
		articles []repository.Article
		response string
	)

	mysql := repository.GlobalMysql
	if articles, err = repository.FindArticleCmd(mysql, repository.Article{}, 0); err != nil {
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = parseToJson(&res, c); err != nil {
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

	argFlg = 1 << 2
	mysql := repository.GlobalMysql
	if articles, err = repository.FindArticleCmd(mysql, argArticle, argFlg); err != nil {
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = parseToJson(&res, c); err != nil {
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

	argFlg = 1 << 4
	mysql := repository.GlobalMysql
	if articles, err = repository.FindArticleCmd(mysql, argArticle, argFlg); err != nil {
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = parseToJson(&res, c); err != nil {
		return
	}

	fmt.Fprint(c.Writer, response)
}

func FindArticleHandlerByCategory(c *gin.Context) {
	var (
		argArticle repository.Article
		articles   []repository.Article
		response   string
		err        error
	)

	// add parameter handling

	mysql := repository.GlobalMysql
	if articles, err = repository.FindArticleCmd(mysql, argArticle, 0); err != nil {
		return
	}

	res := ResponseType{Articles: articles}
	if response, err = parseToJson(&res, c); err != nil {
		return
	}

	fmt.Fprint(c.Writer, response)
}

func FindCategoryHandler(c *gin.Context) {
	var (
		err        error
		categories []repository.Category
		response   string
	)

	mysql := repository.GlobalMysql
	if categories, err = repository.FindCategoryCmd(mysql, repository.Category{}, 0); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	res := CategoryResponseType{Categories: categories}
	if response, err = parseToJson(&res, c); err != nil {
		return
	}

	fmt.Fprint(c.Writer, response)
}
