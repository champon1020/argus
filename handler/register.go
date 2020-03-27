package handler

import (
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/repo"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

// Property of 'categories' has -1 of id absolutely.
// This is because client side cannot fetch categories information interactively,
type RequestArticle struct {
	Id          string          `json:"id"`
	Title       string          `json:"title"`
	Categories  []repo.Category `json:"categories"`
	ContentHash string          `json:"contentHash"`
	ImageHash   string          `json:"imageHash"`
	Private     bool            `json:"isPrivate"`
}

type RequestBody struct {
	Article struct {
		Id          string          `json:"id"`
		Title       string          `json:"title"`
		Categories  []repo.Category `json:"categories"`
		ContentHash string          `json:"contentHash"`
		ImageHash   string          `json:"imageHash"`
		Private     bool            `json:"isPrivate"`
	} `json:"article"`
	MdeContents  string `json:"mdContents"`
	HtmlContents string `json:"htmlContents"`
}

func RegisterArticleController(c *gin.Context) {
	_ = RegisterArticleHandler(c, repo.RegisterArticleCommand)
}

func RegisterArticleHandler(c *gin.Context, repoCmd repo.RegisterArticleCmd) (err error) {
	var body RequestBody

	if err = ParseRequestBody(c.Request, &body); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fp := service.ResolveContentFilePath("articles", body.Article.ContentHash)
	article := repo.Article{
		Title:       body.Article.Title,
		Categories:  body.Article.Categories,
		CreateDate:  time.Now(),
		UpdateDate:  time.Now(),
		ContentHash: filepath.Base(fp),
		ImageHash:   body.Article.ImageHash,
		Private:     body.Article.Private,
	}
	service.GenNewId(service.IdLen, &article.Id)

	htmlFp := fp + "_html"
	mdFp := fp + "_md"

	// output html
	if err = service.OutputFile(htmlFp, []byte(body.HtmlContents)); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// output md
	if err = service.OutputFile(mdFp, []byte(body.MdeContents)); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = repoCmd(*repo.GlobalMysql, article); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		// if err occurred, delete created files.
		_ = service.DeleteFile(htmlFp)
		_ = service.DeleteFile(mdFp)
		return
	}

	_ = service.DeleteFile(
		filepath.Join(
			argus.EnvVars.Get("resource"),
			"drafts",
			body.Article.ContentHash+"_md",
		),
	)

	c.AbortWithStatus(http.StatusOK)
	return
}

func RegisterImageController(c *gin.Context) {
	_ = RegisterImageHandler(c)
}

func RegisterImageHandler(c *gin.Context) (err error) {
	var form *multipart.Form

	if form, err = c.MultipartForm(); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		BasicError.
			SetErr(err).
			SetValues("header", c.Request.Header).
			SetValues("multipart-form", form).
			AppendTo(Errors)
		return
	}

	fileHeaders := form.File["images"]
	path := filepath.Join(argus.EnvVars.Get("resource"), "images")
	if err = service.SaveMultipartFiles(path, fileHeaders); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.Header().Set("location", argus.GlobalConfig.Web.Host+"/manage/images")
	c.AbortWithStatus(http.StatusOK)
	return
}
