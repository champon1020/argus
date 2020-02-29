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
		w    http.ResponseWriter
	)

	w = c.Writer
	if err = ParseRequestBody(&w, c.Request, &body); err != nil {
		return
	}

	mysql := repository.GlobalMysql
	if err = repository.UpdateArticleCmd(mysql, body.Article); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(c.Writer, http.StatusOK)
}
