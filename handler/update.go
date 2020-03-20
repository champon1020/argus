package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/champon1020/argus/repo"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

func UpdateArticleController(c *gin.Context) {
	_ = UpdateArticleHandler(c, repo.UpdateArticleCommand)
}

func UpdateArticleHandler(c *gin.Context, repoCmd repo.UpdateArticleCmd) (err error) {
	var body RequestBody

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

	if err = service.OutputFile(fp, []byte(body.Contents)); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = repoCmd(*repo.GlobalMysql, article); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		service.DeleteFile(fp)
		return
	}

	fmt.Fprint(c.Writer, http.StatusOK)
	return
}

type RequestArticleObjType struct {
	Article repo.Article `json:"article"`
}

func UpdateArticleObjController(c *gin.Context) {
	_ = UpdateArticleObjHandler(c, repo.UpdateArticleCommand)
}

func UpdateArticleObjHandler(c *gin.Context, repoCmd repo.UpdateArticleCmd) (err error) {
	var body RequestBody

	if err = ParseRequestBody(c.Request, &body); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	article := repo.Article{
		Id:          body.Article.Id,
		Title:       body.Article.Title,
		Categories:  body.Article.Categories,
		UpdateDate:  time.Now(),
		ContentHash: body.Article.ContentHash,
		ImageHash:   body.Article.ImageHash,
		Private:     body.Article.Private,
	}

	if err = repoCmd(*repo.GlobalMysql, article); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(c.Writer, http.StatusOK)
	return
}
