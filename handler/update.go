package handler

import (
	"net/http"
	"time"

	"github.com/champon1020/argus/repo"
	"github.com/gin-gonic/gin"
)

func UpdateArticleController(c *gin.Context) {
	_ = UpdateArticleHandler(c, repo.UpdateArticleCommand)
}

func UpdateArticleHandler(c *gin.Context, repoCmd repo.UpdateArticleCmd) (err error) {
	var body RequestBody

	if err = ParseRequestBody(c.Request, &body); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	article := repo.Article{
		Id:         body.Article.Id,
		Title:      body.Article.Title,
		Categories: body.Article.Categories,
		UpdateDate: time.Now(),
		Content:    body.Article.Content,
		ImageHash:  body.Article.ImageHash,
		Private:    body.Article.Private,
	}

	if err = repoCmd(*repo.GlobalMysql, article); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
	return
}
