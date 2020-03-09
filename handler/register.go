package handler

import (
	"fmt"
	"net/http"

	"github.com/champon1020/argus/repository"
	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Article  repository.Article `json:"article"`
	Contents string             `json:"contents"`
}

func RegisterArticleHandler(c *gin.Context) {
	var (
		body RequestBody
		err  error
	)

	if err = ParseRequestBody(c.Request, &body); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	fp := ResolveContentFilePath(body.Article.ContentHash, "articles")
	body.Article.ContentHash = ConvertPathToFileName(fp)
	if err = OutputFile(fp, body.Contents); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	mysql := repository.GlobalMysql
	if err = repository.RegisterArticleCmd(mysql, body.Article); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		DeleteFile(fp)
		return
	}

	fmt.Fprint(c.Writer, http.StatusOK)
}
