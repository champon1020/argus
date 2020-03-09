package handler

import (
	"fmt"
	"net/http"

	"github.com/champon1020/argus/repository"
	"github.com/gin-gonic/gin"
)

func UpdateArticleHandler(c *gin.Context) {
	var (
		body RequestBody
		err  error
	)

	if err = ParseRequestBody(c.Request, &body); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	fp := ResolveContentFilePath(body.Article.ContentUrl, "articles")
	if err = OutputFile(fp, body.Contents); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	mysql := repository.GlobalMysql
	if err = repository.UpdateArticleCmd(mysql, body.Article); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		DeleteFile(fp)
		return
	}
	fmt.Fprint(c.Writer, http.StatusOK)
}
