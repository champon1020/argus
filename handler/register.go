package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/champon1020/argus/service"

	"github.com/champon1020/argus/repository"
	"github.com/gin-gonic/gin"
)

// Property of 'categories' has -1 of id absolutely.
// This is because client side cannot fetch categories information interactively,
type RequestArticle struct {
	Id          int                   `json:"id"`
	Title       string                `json:"title"`
	Categories  []repository.Category `json:"categories"`
	ContentHash string                `json:"contentHash"`
	ImageHash   string                `json:"imageHash"`
	Private     bool                  `json:"private"`
}

type RequestBody struct {
	Article  RequestArticle `json:"article"`
	Contents string         `json:"contents"`
}

func RegisterArticleHandler(c *gin.Context) {
	var (
		body RequestBody
		err  error
	)

	if err = ParseRequestBody(c.Request, &body); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fp := service.ResolveContentFilePath(body.Article.ContentHash, "articles")
	article := repository.Article{
		Id:          body.Article.Id,
		Title:       body.Article.Title,
		Categories:  body.Article.Categories,
		CreateDate:  time.Now(),
		UpdateDate:  time.Now(),
		ContentHash: service.ConvertPathToFileName(fp),
		ImageHash:   body.Article.ImageHash,
		Private:     body.Article.Private,
	}

	if err = service.OutputFile(fp, body.Contents); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mysql := repository.GlobalMysql
	if err = repository.RegisterArticleCmd(mysql, article); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		service.DeleteFile(fp)
		return
	}

	fmt.Fprint(c.Writer, http.StatusOK)
}
