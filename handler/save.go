package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SaveArticleHandler(c *gin.Context) {
	var (
		body RequestBody
		err  error
	)

	if err = ParseRequestBody(c.Request, &body); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	fp := ResolveContentFilePath(body.Article.ContentUrl, "drafts")
	body.Article.ContentUrl = ConvertPathToFileName(fp)
	if err = OutputFile(fp, body.Contents); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		DeleteFile(fp)
		return
	}

	fmt.Fprint(c.Writer, http.StatusOK)
}
