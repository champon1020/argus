package handler

import (
	"fmt"
	"net/http"

	"github.com/champon1020/argus/repository"
	"github.com/gin-gonic/gin"
)

func RegisterArticleHandler(c *gin.Context) {
	body, err := RootHandler(c.Writer, c.Request, "POST")
	if err != nil {
		return
	}

	err = repository.RegisterArticleCmd(mysql, body.Article)
	if err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(c.Writer, http.StatusOK)
}
