package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/champon1020/argus/repository"
)

func FindArticleHandler(c *gin.Context) {
	body, err := RootHandler(c.Writer, c.Request, "GET")
	if err != nil {
		return
	}

	res, err := repository.FindArticleCmd(mysql, body.Article, 0)
	if err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(c.Writer, res)
}

func FindCategoryHandler(c *gin.Context) {
	res, err := repository.FindCategoryCmd(mysql, repository.Category{}, 0)
	if err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(c.Writer, res)
}
