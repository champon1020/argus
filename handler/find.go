package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/repo"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

type ResponseType struct {
	Articles []repo.Article `json:"articles"`
	MaxPage  int            `json:"maxPage"`
}

func FindArticleController(c *gin.Context) {
	_ = FindArticleHandler(c, repo.FindArticleCommand, repo.FindArticlesNumCommand)
}

func FindArticleHandler(
	c *gin.Context,
	repoCmd repo.FindArticleCmd,
	repoNumCmd repo.FindArticleNumCmd,
) (err error) {
	var (
		articles    []repo.Article
		articlesNum int
		response    string
		p           int
	)

	if p, err = ParsePage(c); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ol := ParseOffsetLimit(p)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		argsMask := service.GenMask(repo.Article{}, "Limit")
		if articles, err = repoCmd(*repo.GlobalMysql, repo.Article{}, argsMask, ol); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		numArgsFlg := service.GenMask(repo.Article{})
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, repo.Article{}, numArgsFlg); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	wg.Wait()

	res := ResponseType{
		Articles: articles,
		MaxPage:  GetMaxPage(articlesNum, argus.GlobalConfig.Web.MaxViewArticleNum),
	}
	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
	return
}

type ArticleResponseType struct {
	Article repo.Article `json:"article"`
}

func FindArticleByIdController(c *gin.Context) {
	_ = FindArticleByIdHandler(c, repo.FindArticleCommand)
}

func FindArticleByIdHandler(
	c *gin.Context,
	repoCmd repo.FindArticleCmd,
) (err error) {
	var (
		argArticle repo.Article
		articles   []repo.Article
		response   string
		argsMask   uint32
	)

	if argArticle.Id, err = strconv.Atoi(c.Query("id")); err != nil {
		BasicError.SetErr(err).AppendTo(Errors)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	argsMask = service.GenMask(repo.Article{}, "Id")
	if articles, err = repoCmd(*repo.GlobalMysql, argArticle, argsMask, [2]int{}); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(articles) == 0 {
		err = errors.New("article is not found")
		BasicError.SetErr(err).AppendTo(Errors)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := ArticleResponseType{
		Article: articles[0],
	}
	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
	return
}

func FindArticleByTitleController(c *gin.Context) {
	_ = FindArticleByTitleHandler(c, repo.FindArticleCommand, repo.FindArticlesNumCommand)
}

func FindArticleByTitleHandler(
	c *gin.Context,
	repoCmd repo.FindArticleCmd,
	repoNumCmd repo.FindArticleNumCmd,
) (err error) {
	var (
		argArticle  repo.Article
		articles    []repo.Article
		articlesNum int
		response    string
		argsMask    uint32
		p           int
	)

	argArticle.Title = c.Query("title")

	if p, err = ParsePage(c); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ol := ParseOffsetLimit(p)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		argsMask = service.GenMask(repo.Article{}, "Title", "Limit")
		if articles, err = repoCmd(*repo.GlobalMysql, argArticle, argsMask, ol); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		numArgsFlg := service.GenMask(repo.Article{})
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, repo.Article{}, numArgsFlg); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	wg.Wait()

	res := ResponseType{
		Articles: articles,
		MaxPage:  GetMaxPage(articlesNum, argus.GlobalConfig.Web.MaxViewArticleNum),
	}
	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
	return
}

func FindArticleByCreateDateController(c *gin.Context) {
	_ = FindArticleByCreateDateHandler(c, repo.FindArticleCommand, repo.FindArticlesNumCommand)
}

func FindArticleByCreateDateHandler(
	c *gin.Context,
	repoCmd repo.FindArticleCmd,
	repoNumCmd repo.FindArticleNumCmd,
) (err error) {
	var (
		argArticle  repo.Article
		articles    []repo.Article
		articlesNum int
		response    string
		argsMask    uint32
		p           int
	)

	if argArticle.CreateDate, err = time.Parse(time.RFC3339, c.Query("createDate")); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		TimeParseError.SetErr(err).AppendTo(Errors)
		return
	}

	if p, err = ParsePage(c); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ol := ParseOffsetLimit(p)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		argsMask = service.GenMask(repo.Article{}, "CreateDate", "Limit")
		if articles, err = repoCmd(*repo.GlobalMysql, argArticle, argsMask, ol); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		numArgsFlg := service.GenMask(repo.Article{})
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, repo.Article{}, numArgsFlg); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	wg.Wait()

	res := ResponseType{
		Articles: articles,
		MaxPage:  GetMaxPage(articlesNum, argus.GlobalConfig.Web.MaxViewArticleNum),
	}
	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
	return
}

func FindArticleByCategoryController(c *gin.Context) {
	_ = FindArticleByCategoryHandler(
		c,
		repo.FindArticleByCategoryCommand,
		repo.FindArticlesNumByCategoryCommand)
}

func FindArticleByCategoryHandler(
	c *gin.Context,
	repoCmd repo.FindArticleByCategoryCmd,
	repoNumCmd repo.FindArticlesNumByCategoryCmd,
) (err error) {
	var (
		articles    []repo.Article
		articlesNum int
		response    string
		argsMask    uint32
		p           int
	)

	categoryNames := c.QueryArray("category")

	if p, err = ParsePage(c); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ol := ParseOffsetLimit(p)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		argsMask = service.GenMask(repo.Article{}, "Limit")
		if articles, err = repoCmd(*repo.GlobalMysql, categoryNames, argsMask, ol); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, categoryNames); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	wg.Wait()

	res := ResponseType{
		Articles: articles,
		MaxPage:  GetMaxPage(articlesNum, argus.GlobalConfig.Web.MaxViewArticleNum),
	}
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
		argsMask   uint32
	)

	argsMask = service.GenMask(repo.Category{})
	if categories, err = repoCmd(*repo.GlobalMysql, repo.Category{}, argsMask, [2]int{}); err != nil {
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
	Drafts    []repo.Draft `json:"drafts"`
	DraftsNum int          `json:"draftsNum"`
	MaxPage   int          `json:"maxPage"`
}

func FindDraftController(c *gin.Context) {
	_ = FindDraftHandler(c, repo.FindDraftCommand, repo.FindDraftsNumCommand)
}

func FindDraftHandler(
	c *gin.Context,
	repoCmd repo.FindDraftCmd,
	repoNumCmd repo.FindDraftNumCmd,
) (err error) {
	var (
		drafts    []repo.Draft
		draftsNum int
		response  string
		argsMask  uint32
		p         int
	)

	if p, err = ParsePage(c); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ol := ParseOffsetLimit(p)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		argsMask = service.GenMask(repo.Draft{}, "Limit")
		if drafts, err = repoCmd(*repo.GlobalMysql, repo.Draft{}, argsMask, ol); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		numArgsFlg := service.GenMask(repo.Article{})
		if draftsNum, err = repoNumCmd(*repo.GlobalMysql, repo.Draft{}, numArgsFlg); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	wg.Wait()

	maxPage := GetMaxPage(draftsNum, argus.GlobalConfig.Web.MaxViewSettingArticleNum)
	res := DraftResponseType{
		Drafts:    drafts,
		DraftsNum: draftsNum,
		MaxPage:   maxPage,
	}
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
		p         int
	)

	dirPath := filepath.Join(argus.EnvVars.Get("resource"), "images")
	if files, err = service.GetFileList(dirPath); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if p, err = ParsePage(c); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	offset := (p - 1) * argus.GlobalConfig.Web.MaxViewImageNum

	if offset >= len(files) {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		err = errors.New("error happened")
		BasicError.
			SetValues("len(files)", len(files)).
			SetValues("offset", offset).
			SetErr(err).
			AppendTo(Errors)
		return
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
