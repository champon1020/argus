package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/champon1020/argus/repository"
	"github.com/gin-gonic/gin"
)

func UpdateArticleHandler(c *gin.Context) {
	var (
		body RequestBody
		err  error
	)

	if err = ParseRequestBody(c.Request, &body); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fp := ResolveContentFilePath(body.Article.ContentHash, "articles")
	article := repository.Article{
		Id:          body.Article.Id,
		Title:       body.Article.Title,
		Categories:  body.Article.Categories,
		UpdateDate:  time.Now(),
		ContentHash: body.Article.ContentHash,
		ImageHash:   body.Article.ImageHash,
		Private:     body.Article.Private,
	}

	if err = OutputFile(fp, body.Contents); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mysql := repository.GlobalMysql
	if err = repository.UpdateArticleCmd(mysql, article); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		DeleteFile(fp)
		return
	}
	fmt.Fprint(c.Writer, http.StatusOK)
}
