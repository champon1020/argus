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

	if body, err = RootHandler(c.Writer, c.Request, "PUT"); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	if err = repository.UpdateArticleCmd(mysql, body.Article); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(c.Writer, http.StatusOK)
}
