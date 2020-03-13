package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	repo "github.com/champon1020/argus/repository"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

type ResponseType struct {
	Articles []repo.Article `json:"articles"`
}

func FindArticleController(c *gin.Context) {
	FindArticleHandler(c, repo.FindArticleCommand)
}

func FindArticleHandler(c *gin.Context, repoCmd repo.FindArticleCmd) {
	var (
		articles []repo.Article
		response string
		argFlg   uint32
		err      error
	)

	argFlg = service.GenFlg(repo.Article{}, "Limit")
	if articles, err = repoCmd(*repo.GlobalMysql, repo.Article{}, argFlg); err != nil {
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

func FindArticleByIdController(c *gin.Context) {
	FindArticleByIdHandler(c, repo.FindArticleCommand)
}

func FindArticleByIdHandler(c *gin.Context, repoCmd repo.FindArticleCmd) {
	var (
		argArticle repo.Article
		articles   []repo.Article
		response   string
		argFlg     uint32
		err        error
	)

	if argArticle.Id, err = strconv.Atoi(c.Query("id")); err != nil {
		BasicError.SetErr(err).AppendTo(&Errors)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	argFlg = service.GenFlg(repo.Article{}, "Id", "Limit")
	if articles, err = repoCmd(*repo.GlobalMysql, argArticle, argFlg); err != nil {
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

func FindArticleByTitleController(c *gin.Context) {
	FindArticleByTitleHandler(c, repo.FindArticleCommand)
}

func FindArticleByTitleHandler(c *gin.Context, repoCmd repo.FindArticleCmd) {
	var (
		argArticle repo.Article
		articles   []repo.Article
		response   string
		argFlg     uint32
		err        error
	)

	argArticle.Title = c.Query("title")

	argFlg = service.GenFlg(repo.Article{}, "Title", "Limit")
	if articles, err = repoCmd(*repo.GlobalMysql, argArticle, argFlg); err != nil {
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

func FindArticleByCreateDateController(c *gin.Context) {
	FindArticleByCreateDateHandler(c, repo.FindArticleCommand)
}

func FindArticleByCreateDateHandler(c *gin.Context, repoCmd repo.FindArticleCmd) {
	var (
		argArticle repo.Article
		articles   []repo.Article
		response   string
		argFlg     uint32
		err        error
	)

	if argArticle.CreateDate, err = time.Parse(time.RFC3339, c.Query("createDate")); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		TimeParseError.SetErr(err).AppendTo(&Errors)
		return
	}

	argFlg = service.GenFlg(repo.Article{}, "CreateDate", "Limit")
	if articles, err = repoCmd(*repo.GlobalMysql, argArticle, argFlg); err != nil {
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

func FindArticleByCategoryController(c *gin.Context) {
	FindArticleByCategoryHandler(c, repo.FindArticleByCategoryCommand)
}

func FindArticleByCategoryHandler(c *gin.Context, repoCmd repo.FindArticleByCategoryCmd) {
	var (
		articles []repo.Article
		response string
		argFlg   uint32
		err      error
	)

	categoryNames := c.QueryArray("category")
	argFlg = service.GenFlg(repo.Article{}, "Limit")
	if articles, err = repoCmd(*repo.GlobalMysql, categoryNames, argFlg); err != nil {
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
	Categories []repo.CategoryResponse `json:"categories"`
}

func FindCategoryController(c *gin.Context) {
	FindCategoryHandler(c, repo.FindCategoryCommand)
}

func FindCategoryHandler(c *gin.Context, repoCmd repo.FindCategoryCmd) {
	var (
		categories []repo.CategoryResponse
		response   string
		argFlg     uint32
		err        error
	)

	argFlg = service.GenFlg(repo.Category{}, "Limit")
	if categories, err = repoCmd(*repo.GlobalMysql, repo.Category{}, argFlg); err != nil {
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
	Drafts []repo.Draft `json:"drafts"`
}

func FindDraftController(c *gin.Context) {
	FindDraftHandler(c, repo.FindDraftCommand)
}

func FindDraftHandler(c *gin.Context, repoCmd repo.FindDraftCmd) {
	var (
		drafts   []repo.Draft
		response string
		argFlg   uint32
		err      error
	)

	argFlg = service.GenFlg(repo.Draft{}, "Limit")
	if drafts, err = repoCmd(*repo.GlobalMysql, repo.Draft{}, argFlg); err != nil {
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
