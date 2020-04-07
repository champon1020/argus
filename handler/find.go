package handler

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/repo"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

// Find article by id
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
	id := c.Query("id")
	res := new(ArticleResponseType)

	var articles []repo.Article
	option := &service.QueryOption{
		Args: []*service.QueryArgs{
			{
				Value: id,
				Name:  "Id",
				Ope:   service.Eq,
			},
		},
	}
	if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(articles) > 0 {
		res.Article = articles[0]
	}

	c.JSON(http.StatusOK, res)
	return
}

type ArticleSortedResponseType struct {
	Article repo.Article `json:"article"`
	Next    repo.Article `json:"next"`
	Prev    repo.Article `json:"prev"`
}

func FindArticleBySortedIdController(c *gin.Context) {
	_ = FindArticleBySortedIdHandler(c, repo.FindArticleCommand)
}

func FindArticleBySortedIdHandler(
	c *gin.Context,
	repoCmd repo.FindArticleCmd,
) (err error) {
	sortedId := c.Query("sortedId")
	res := new(ArticleSortedResponseType)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	// get current and previous article
	go func() {
		defer wg.Done()
		var articles []repo.Article
		option := &service.QueryOption{
			Args: []*service.QueryArgs{
				{
					Value: false,
					Name:  "Private",
					Ope:   service.Eq,
				},
				{
					Value: sortedId,
					Name:  "SortedId",
					Ope:   service.Ge,
				},
			},
			Limit:  2,
			Offset: 0,
			Order:  "sorted_id",
		}
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if len(articles) > 0 {
			res.Article = articles[0]
		}
		if len(articles) > 1 {
			res.Prev = articles[1]
		}
	}()
	// get next article
	go func() {
		defer wg.Done()
		var articles []repo.Article
		option := &service.QueryOption{
			Args: []*service.QueryArgs{
				{
					Value: false,
					Name:  "Private",
					Ope:   service.Eq,
				},
				{
					Value: sortedId,
					Name:  "SortedId",
					Ope:   service.Lt,
				},
			},
			Limit:  1,
			Offset: 0,
			Order:  "sorted_id",
			Desc:   true,
		}
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if len(articles) == 1 {
			res.Next = articles[0]
		}
	}()
	wg.Wait()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, res)
	return
}

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
		p           int
	)

	if p, err = ParsePage(c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		ol := ParseOffsetLimit(p, argus.GlobalConfig.Web.MaxViewArticleNum)
		option := &service.QueryOption{
			Args: []*service.QueryArgs{
				{
					Value: false,
					Name:  "Private",
					Ope:   service.Eq,
				},
			},
			Limit:  ol[1],
			Offset: ol[0],
			Order:  "sorted_id",
			Desc:   true,
		}
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, service.DefaultOption); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()
	wg.Wait()
	if err != nil {
		return
	}

	res := ResponseType{
		Articles: articles,
		MaxPage:  GetMaxPage(articlesNum, argus.GlobalConfig.Web.MaxViewArticleNum),
	}

	c.JSON(http.StatusOK, res)
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
		p           int
	)

	if p, err = ParsePage(c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		ol := ParseOffsetLimit(p, argus.GlobalConfig.Web.MaxViewArticleNum)
		option := &service.QueryOption{
			Args: []*service.QueryArgs{
				{
					Value: "%" + c.Query("title") + "%",
					Name:  "Title",
					Ope:   service.Like,
				},
				{
					Value: false,
					Name:  "Private",
					Ope:   service.Eq,
				},
			},
			Limit:  ol[1],
			Offset: ol[0],
			Order:  "sorted_id",
			Desc:   true,
		}
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		option := &service.QueryOption{
			Args: []*service.QueryArgs{
				{
					Value: "%" + c.Query("title") + "%",
					Name:  "Title",
					Ope:   service.Like,
				},
			},
		}
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()
	wg.Wait()
	if err != nil {
		return
	}

	res := ResponseType{
		Articles: articles,
		MaxPage:  GetMaxPage(articlesNum, argus.GlobalConfig.Web.MaxViewArticleNum),
	}

	c.JSON(http.StatusOK, res)
	return
}

func FindArticleByCreateDateController(c *gin.Context) {
	_ = FindArticleByCreateDateHandler(
		c,
		repo.FindArticleCommand,
		repo.FindArticlesNumCommand)
}

func FindArticleByCreateDateHandler(
	c *gin.Context,
	repoCmd repo.FindArticleCmd,
	repoNumCmd repo.FindArticleNumCmd,
) (err error) {
	var (
		articles    []repo.Article
		articlesNum int
		p           int
	)

	if p, err = ParsePage(c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		ol := ParseOffsetLimit(p, argus.GlobalConfig.Web.MaxViewArticleNum)
		option := &service.QueryOption{
			Args: []*service.QueryArgs{
				{
					Value: c.Query("CreateDate"),
					Name:  "CreateDate",
					Ope:   service.Eq,
				},
				{
					Value: false,
					Name:  "Private",
					Ope:   service.Eq,
				},
			},
			Limit:  ol[1],
			Offset: ol[0],
			Order:  "sorted_id",
			Desc:   true,
		}
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		option := &service.QueryOption{
			Args: []*service.QueryArgs{
				{
					Value: c.Query("CreateDate"),
					Name:  "CreateDate",
					Ope:   service.Eq,
				},
			},
		}
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()
	wg.Wait()
	if err != nil {
		return
	}

	res := ResponseType{
		Articles: articles,
		MaxPage:  GetMaxPage(articlesNum, argus.GlobalConfig.Web.MaxViewArticleNum),
	}

	c.JSON(http.StatusOK, res)
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
		p           int
	)

	if p, err = ParsePage(c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	categoryNames := c.QueryArray("category")

	// get articles
	go func() {
		defer wg.Done()
		ol := ParseOffsetLimit(p, argus.GlobalConfig.Web.MaxViewArticleNum)
		option := &service.QueryOption{
			Args: []*service.QueryArgs{
				{
					Value: false,
					Name:  "Private",
					Ope:   service.Eq,
				},
			},
			Limit:  ol[1],
			Offset: ol[0],
			Order:  "sorted_id",
			Desc:   true,
		}
		for _, c := range categoryNames {
			option.Args = append(option.Args, &service.QueryArgs{
				Value: c,
				Name:  "Name",
				Ope:   service.Eq,
			})
		}
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		option := &service.QueryOption{
			Args: []*service.QueryArgs{
				{
					Value: false,
					Name:  "Private",
					Ope:   service.Eq,
				},
			},
		}
		for _, c := range categoryNames {
			option.Args = append(option.Args, &service.QueryArgs{
				Value: c,
				Name:  "Name",
				Ope:   service.Eq,
			})
		}
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()
	wg.Wait()
	if err != nil {
		return
	}

	res := ResponseType{
		Articles: articles,
		MaxPage:  GetMaxPage(articlesNum, argus.GlobalConfig.Web.MaxViewArticleNum),
	}

	c.JSON(http.StatusOK, res)
	return
}

func FindAllArticleController(c *gin.Context) {
	_ = FindAllArticleHandler(
		c,
		repo.FindArticleCommand,
		repo.FindArticlesNumCommand)
}

func FindAllArticleHandler(
	c *gin.Context,
	repoCmd repo.FindArticleCmd,
	repoNumCmd repo.FindArticleNumCmd,
) (err error) {
	var (
		articles    []repo.Article
		articlesNum int
		p           int
	)

	if p, err = ParsePage(c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := new(ResponseType)
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		ol := ParseOffsetLimit(p, argus.GlobalConfig.Web.MaxViewSettingArticleNum)
		option := &service.QueryOption{
			Args:   []*service.QueryArgs{},
			Limit:  ol[1],
			Offset: ol[0],
			Order:  "sorted_id",
			Desc:   true,
		}
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		res.Articles = articles
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		if articlesNum, err = repoNumCmd(*repo.GlobalMysql, service.DefaultOption); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		res.MaxPage = GetMaxPage(articlesNum, argus.GlobalConfig.Web.MaxViewSettingArticleNum)
	}()
	wg.Wait()

	c.JSON(http.StatusOK, res)
	return
}

type PickUpArticleResponse struct {
	Articles []repo.Article `json:"articles"`
}

func FindPickUpArticleController(c *gin.Context) {
	_ = FindPickUpArticleHandler(c, repo.FindArticleCommand)
}

func FindPickUpArticleHandler(c *gin.Context, repoCmd repo.FindArticleCmd) (err error) {
	res := new(PickUpArticleResponse)
	pickupId := argus.GlobalConfig.Web.Pickup

	for _, id := range pickupId {
		var articles []repo.Article
		option := &service.QueryOption{
			Args: []*service.QueryArgs{
				{
					Value: id,
					Name:  "SortedId",
					Ope:   service.Eq,
				},
			},
		}
		if articles, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		res.Articles = append(res.Articles, articles...)
	}

	c.JSON(http.StatusOK, res)
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
	var categories []repo.CategoryResponse

	if categories, err = repoCmd(*repo.GlobalMysql, service.DefaultOption); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := CategoryResponseType{Categories: categories}

	c.JSON(http.StatusOK, res)
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
		p         int
	)

	if p, err = ParsePage(c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// get articles
	go func() {
		defer wg.Done()
		ol := ParseOffsetLimit(p, argus.GlobalConfig.Web.MaxViewArticleNum)
		option := &service.QueryOption{
			Limit:  ol[1],
			Offset: ol[0],
			Order:  "sorted_id",
			Desc:   true,
		}
		if drafts, err = repoCmd(*repo.GlobalMysql, option); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()
	// get the number of total articles
	go func() {
		defer wg.Done()
		if draftsNum, err = repoNumCmd(*repo.GlobalMysql, service.DefaultOption); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}()
	wg.Wait()
	if err != nil {
		return
	}

	maxPage := GetMaxPage(draftsNum, argus.GlobalConfig.Web.MaxViewSettingArticleNum)
	res := DraftResponseType{
		Drafts:    drafts,
		DraftsNum: draftsNum,
		MaxPage:   maxPage,
	}

	c.JSON(http.StatusOK, res)
	return
}

type DraftObjectResponse struct {
	Draft repo.Draft `json:"draft"`
}

func FindDraftByIdController(c *gin.Context) {
	_ = FindDraftByIdHandler(c, repo.FindDraftCommand)
}

func FindDraftByIdHandler(c *gin.Context, repoCmd repo.FindDraftCmd) (err error) {
	var drafts []repo.Draft

	id := c.Query("id")

	res := new(DraftObjectResponse)
	option := &service.QueryOption{
		Args: []*service.QueryArgs{
			{
				Value: id,
				Name:  "id",
				Ope:   service.Eq,
			},
		},
	}

	if drafts, err = repoCmd(*repo.GlobalMysql, option); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(drafts) > 0 {
		res.Draft = drafts[0]
	}

	c.JSON(http.StatusOK, res)
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
		files     []os.FileInfo
		fileNames []string
		p         int
	)

	dirPath := filepath.Join(argus.EnvVars.Get("resource"), "images")
	if files, err = service.GetFileList(dirPath); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if p, err = ParsePage(c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	mx := argus.GlobalConfig.Web.MaxViewImageNum
	offset := (p - 1) * mx

	if offset >= len(files) {
		c.AbortWithStatus(http.StatusInternalServerError)
		err = errors.New("error happened")
		BasicError.
			SetValues("len(files)", len(files)).
			SetValues("offset", offset).
			SetErr(err).
			AppendTo(Errors)
		return
	}

	res.Next = true
	for i := offset; i < offset+mx; i++ {
		if len(files) <= i {
			res.Next = false
			break
		}
		fileNames = append(fileNames, files[i].Name())
	}

	res.Images = fileNames

	c.JSON(http.StatusOK, res)
	return
}
