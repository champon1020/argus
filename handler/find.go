package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/champon1020/argus/repository"
)

func FindArticleHandler(c *gin.Context) {
	var (
		body RequestBody
		err  error
		res  []repository.Article
	)

	if body, err = RootHandler(c.Writer, c.Request, "GET"); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	if res, err = repository.FindArticleCmd(mysql, body.Article, 0); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, res)
}

func FindCategoryHandler(c *gin.Context) {
	var (
		err error
		res []repository.Category
	)

	if res, err = repository.FindCategoryCmd(mysql, repository.Category{}, 0); err != nil {
		fmt.Fprint(c.Writer, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, res)
}
