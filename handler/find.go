package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

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

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		ol := ParseOffsetLimit(p)
		option := &service.QueryOption{
			Args:   []interface{}{false},
			Aom:    map[string]service.Ope{"Private": service.Eq},
			Limit:  ol[1],
			Offset: ol[0],
			Order:  "sorted_id",
			Desc:   true,
		}
		option.BuildArgs()
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, service.DefaultOption); err != nil {
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
	Next    repo.Article `json:"next"`
	Prev    repo.Article `json:"prev"`
}

func FindArticleByIdController(c *gin.Context) {
	_ = FindArticleByIdHandler(c, repo.FindArticleCommand)
}

func FindArticleByIdHandler(
	c *gin.Context,
	repoCmd repo.FindArticleCmd,
) (err error) {
	var (
		articles []repo.Article
		response string
	)

	sortedId := c.Query("sortedId")
	res := new(ArticleResponseType)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	// get current and previous article
	go func() {
		defer wg.Done()
		option := &service.QueryOption{
			Args: []interface{}{sortedId, false},
			Aom: map[string]service.Ope{
				"SortedId": service.Ge,
				"Private":  service.Eq,
			},
			Limit:  2,
			Offset: 0,
			Order:  "sorted_id",
		}
		option.BuildArgs()
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Article = articles[0]
		if len(articles) == 2 {
			res.Prev = articles[1]
		}
	}()
	// get next article
	go func() {
		defer wg.Done()
		option := &service.QueryOption{
			Args: []interface{}{sortedId, false},
			Aom: map[string]service.Ope{
				"SortedId": service.Lt,
				"Private":  service.Eq,
			},
			Limit:  1,
			Offset: 0,
			Order:  "sorted_id",
			Desc:   true,
		}
		option.BuildArgs()
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(articles) == 1 {
			res.Next = articles[0]
		}
	}()
	wg.Wait()

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
		articles    []repo.Article
		articlesNum int
		response    string
		p           int
	)

	if p, err = ParsePage(c); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		ol := ParseOffsetLimit(p)
		option := &service.QueryOption{
			Args: []interface{}{c.Query("title"), false},
			Aom: map[string]service.Ope{
				"Title":   service.Eq,
				"Private": service.Eq,
			},
			Limit:  ol[1],
			Offset: ol[0],
			Order:  "sorted_id",
			Desc:   true,
		}
		option.BuildArgs()
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		option := &service.QueryOption{
			Args: []interface{}{c.Query("title")},
			Aom:  map[string]service.Ope{"Title": service.Eq},
		}
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, option); err != nil {
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
		articles    []repo.Article
		articlesNum int
		response    string
		p           int
	)

	createDate := c.Query("createDate")

	if p, err = ParsePage(c); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		ol := ParseOffsetLimit(p)
		option := &service.QueryOption{
			Args: []interface{}{createDate, false},
			Aom: map[string]service.Ope{
				"CreateDate": service.Eq,
				"Private":    service.Eq,
			},
			Limit:  ol[1],
			Offset: ol[0],
			Order:  "sorted_id",
			Desc:   true,
		}
		option.BuildArgs()
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		option := &service.QueryOption{
			Args: []interface{}{createDate},
			Aom:  map[string]service.Ope{"CreateDate": service.Eq},
		}
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, option); err != nil {
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
		p           int
	)

	if p, err = ParsePage(c); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	categoryNames := c.QueryArray("category")

	// get articles
	go func() {
		defer wg.Done()
		ol := ParseOffsetLimit(p)
		option := &service.QueryOption{
			Args:   []interface{}{false},
			Aom:    map[string]service.Ope{"Private": service.Eq},
			Limit:  ol[1],
			Offset: ol[0],
			Order:  "sorted_id",
			Desc:   true,
		}
		option.BuildArgs()
		if articles, err = repoCmd(*repo.GlobalMysql, categoryNames, option); err != nil {
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
	)

	if categories, err = repoCmd(*repo.GlobalMysql, service.DefaultOption); err != nil {
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
		p         int
	)

	if p, err = ParsePage(c); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		ol := ParseOffsetLimit(p)
		option := &service.QueryOption{
			Limit:  ol[1],
			Offset: ol[0],
			Order:  "sorted_id",
			Desc:   true,
		}
		option.BuildArgs()
		if drafts, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		if draftsNum, err = repoNumCmd(*repo.GlobalMysql, service.DefaultOption); err != nil {
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
