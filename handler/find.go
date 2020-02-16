package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/champon1020/argus/repository"
)

func FindArticleHandler(c *gin.Context) {
	body, err := RootHandler(c.Writer, c.Request, "GET", "application/json")
	if err != nil {
		return
	}

	res := repository.FindArticleCmd(mysql, body.Article, 0)
	fmt.Fprint(c.Writer, res)
}

func FindCategoryHandler(w http.ResponseWriter, r *http.Request) {
	res := repository.FindCategoryCmd(mysql, repository.Category{}, 0)
	fmt.Fprint(w, res)
}
