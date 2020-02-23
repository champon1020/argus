package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/champon1020/argus/repository"
)

type QueryParam struct {
	Article repository.Article
}

type ResponseType struct {
	Articles []repository.Article `json:"articles"`
}

func parseQueryParam(param *QueryParam, c *gin.Context) (err error) {
	param.Article.Id, err = strconv.Atoi(c.Query("id"))
	param.Article.Title = c.Query("title")
	param.Article.CreateDate, err = time.Parse(time.RFC3339, c.Query("create_date"))
	param.Article.UpdateDate, err = time.Parse(time.RFC3339, c.Query("update_date"))
	param.Article.ContentUrl = c.Query("content_url")
	param.Article.ImageUrl = c.Query("image_url")
	param.Article.Private, err = strconv.ParseBool(c.Query("private"))
	return
}

func FindArticleHandler(c *gin.Context) {
	var (
		err error
		res []repository.Article
	)

	if res, err = repository.FindArticleCmd(mysql, repository.Article{}, 0); err != nil {
		return
	}

	fmt.Fprint(c.Writer, res)
}

func FindArticleHandlerByTitle(c *gin.Context) {
	var (
		argArticle repository.Article
		err        error
		res        []repository.Article
		argFlg     uint32
	)

	argArticle.Title = c.Query("title")

	argFlg = 1 << 2
	if res, err = repository.FindArticleCmd(mysql, argArticle, argFlg); err != nil {
		return
	}

	fmt.Fprint(c.Writer, res)
}

func FindArticleHandlerByCreateDate(c *gin.Context) {
	var (
		argArticle repository.Article
		err        error
		res        []repository.Article
		argFlg     uint32
	)

	argFlg = 1 << 4
	if res, err = repository.FindArticleCmd(mysql, argArticle, argFlg); err != nil {
		return
	}

	fmt.Fprint(c.Writer, res)
}

func FindArticleHandlerByCategory(c *gin.Context) {
	var (
		argArticle repository.Article
		err        error
		res        []repository.Article
	)

	// add parameter handling

	if res, err = repository.FindArticleCmd(mysql, argArticle, 0); err != nil {
		return
	}

	fmt.Fprint(c.Writer, res)
}

func FindCategoryHandler(c *gin.Context) {
	var (
		err error
		res []repository.Category
	)

	if res, err = repository.FindCategoryCmd(mysql, repository.Category{}, 0); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, res)
}
