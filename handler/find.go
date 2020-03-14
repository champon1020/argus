package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/champon1020/argus"

	repo "github.com/champon1020/argus/repository"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

type ResponseType struct {
	Articles []repo.Article `json:"articles"`
}

func FindArticleController(c *gin.Context) {
	_ = FindArticleHandler(c, repo.FindArticleCommand)
}

func FindArticleHandler(c *gin.Context, repoCmd repo.FindArticleCmd) (err error) {
	var (
		articles []repo.Article
		response string
		argFlg   uint32
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
	return
}

func FindArticleByIdController(c *gin.Context) {
	_ = FindArticleByIdHandler(c, repo.FindArticleCommand)
}

func FindArticleByIdHandler(c *gin.Context, repoCmd repo.FindArticleCmd) (err error) {
	var (
		argArticle repo.Article
		articles   []repo.Article
		response   string
		argFlg     uint32
	)

	if argArticle.Id, err = strconv.Atoi(c.Query("id")); err != nil {
		BasicError.SetErr(err).AppendTo(Errors)
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
	return
}

func FindArticleByTitleController(c *gin.Context) {
	_ = FindArticleByTitleHandler(c, repo.FindArticleCommand)
}

func FindArticleByTitleHandler(c *gin.Context, repoCmd repo.FindArticleCmd) (err error) {
	var (
		argArticle repo.Article
		articles   []repo.Article
		response   string
		argFlg     uint32
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
	return
}

func FindArticleByCreateDateController(c *gin.Context) {
	_ = FindArticleByCreateDateHandler(c, repo.FindArticleCommand)
}

func FindArticleByCreateDateHandler(c *gin.Context, repoCmd repo.FindArticleCmd) (err error) {
	var (
		argArticle repo.Article
		articles   []repo.Article
		response   string
		argFlg     uint32
	)

	if argArticle.CreateDate, err = time.Parse(time.RFC3339, c.Query("createDate")); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		TimeParseError.SetErr(err).AppendTo(Errors)
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
	return
}

func FindArticleByCategoryController(c *gin.Context) {
	_ = FindArticleByCategoryHandler(c, repo.FindArticleByCategoryCommand)
}

func FindArticleByCategoryHandler(c *gin.Context, repoCmd repo.FindArticleByCategoryCmd) (err error) {
	var (
		articles []repo.Article
		response string
		argFlg   uint32
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
	return
}

// Response type of Category
type CategoryResponseType struct {
	Categories []repo.CategoryResponse `json:"categories"`
}

func FindCategoryController(c *gin.Context) {
	_ = FindCategoryHandler(c, repo.FindCategoryCommand)
}

func FindCategoryHandler(c *gin.Context, repoCmd repo.FindCategoryCmd) (err error) {
	var (
		categories []repo.CategoryResponse
		response   string
		argFlg     uint32
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
	return
}

// Response type of Draft
type DraftResponseType struct {
	Drafts []repo.Draft `json:"drafts"`
}

func FindDraftController(c *gin.Context) {
	_ = FindDraftHandler(c, repo.FindDraftCommand)
}

func FindDraftHandler(c *gin.Context, repoCmd repo.FindDraftCmd) (err error) {
	var (
		drafts   []repo.Draft
		response string
		argFlg   uint32
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
	return
}

// Response type of Image
type ImageResponseType struct {
	Images []string `json:"images"`
	Next   bool     `json:"next"`
}

func FindImageController(c *gin.Context) {
	_ = FindImageHandler(c)
}

func FindImageHandler(c *gin.Context) (err error) {
	var (
		res       ImageResponseType
		response  string
		files     []os.FileInfo
		fileNames []string
	)

	// get page
	p := 1
	if pp, ok := c.GetQuery("p"); ok {
		if p, err = strconv.Atoi(pp); err != nil {
			BasicError.SetErr(err).AppendTo(Errors)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	dirPath := filepath.Join(argus.EnvVars.Get("resource"), "images")
	if files, err = service.GetFileList(dirPath); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	offset := (p - 1) * argus.GlobalConfig.Web.MaxViewImageNum
	if offset >= len(files) {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		argus.Logger.Println("Exceed the number of existing files")
		return errors.New("error happened")
	}

	for i := offset; i < len(files); i++ {
		if i >= argus.GlobalConfig.Web.MaxViewImageNum {
			res.Next = true
			break
		}
		fileNames = append(fileNames, files[i].Name())
	}

	res.Images = fileNames

	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
	return
}
