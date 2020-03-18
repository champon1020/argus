package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
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
	Private     bool            `json:"private"`
}

type RequestBody struct {
	Article  RequestArticle `json:"article"`
	Contents string         `json:"contents"`
}

func RegisterArticleController(c *gin.Context) {
	_ = RegisterArticleHandler(c, repo.RegisterArticleCommand)
}

func RegisterArticleHandler(c *gin.Context, repoCmd repo.RegisterArticleCmd) (err error) {
	var body RequestBody

	if err = ParseRequestBody(c.Request, &body); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fp := service.ResolveContentFilePath(body.Article.ContentHash, "articles")
	article := repo.Article{
		Title:       body.Article.Title,
		Categories:  body.Article.Categories,
		CreateDate:  time.Now(),
		UpdateDate:  time.Now(),
		ContentHash: service.ConvertPathToFileName(fp),
		ImageHash:   body.Article.ImageHash,
		Private:     body.Article.Private,
	}
	service.GenNewId(service.IdLen, &article.Id)

	if err = service.OutputFile(fp, []byte(body.Contents)); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = repoCmd(*repo.GlobalMysql, article); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		_ = service.DeleteFile(fp)
		return
	}

	fmt.Fprint(c.Writer, http.StatusOK)
	return
}

func RegisterImageController(c *gin.Context) {
	_ = RegisterImageHandler(c)
}

func RegisterImageHandler(c *gin.Context) (err error) {
	var form *multipart.Form

	if form, err = c.MultipartForm(); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		BasicError.
			SetErr(err).
			SetValues("header", c.Request.Header).
			SetValues("multipart-form", form).
			AppendTo(Errors)
		return
	}

	fileHeaders := form.File["images"]
	path := argus.EnvVars.Get("resource")
	if err = service.SaveMultipartFiles(path, fileHeaders); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, http.StatusOK)
	return
}
