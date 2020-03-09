package handler

import (
	"fmt"
	"net/http"
	"time"

	repo "github.com/champon1020/argus/repository"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

func UpdateArticleController(c *gin.Context) {
	UpdateArticleHandler(c, repo.UpdateArticleCommand)
}

func UpdateArticleHandler(c *gin.Context, repoCmd repo.UpdateArticleCmd) {
	var (
		body RequestBody
		err  error
	)

	if err = ParseRequestBody(c.Request, &body); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fp := service.ResolveContentFilePath(body.Article.ContentHash, "articles")
	article := repo.Article{
		Id:          body.Article.Id,
		Title:       body.Article.Title,
		Categories:  body.Article.Categories,
		UpdateDate:  time.Now(),
		ContentHash: body.Article.ContentHash,
		ImageHash:   body.Article.ImageHash,
		Private:     body.Article.Private,
	}

	if err = service.OutputFile(fp, body.Contents); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mysql := repo.GlobalMysql
	if err = repoCmd(mysql, article); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		service.DeleteFile(fp)
		return
	}
	fmt.Fprint(c.Writer, http.StatusOK)
}
