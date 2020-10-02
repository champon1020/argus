package handler

import (
	"net/http"

	"github.com/champon1020/argus/v2/model"
	"github.com/gin-gonic/gin"
)

// FindArticlesList gets all public articles.
func FindArticlesList(c *gin.Context, db model.DatabaseIface) (err error) {
	var (
		a   []model.Article
		num int
		p   int
	)

	if p, err = ParsePage(c); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if num, err = ParseNum(c); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = db.FindArticlesList(&a, &model.QueryOptions{
		Limit:   num,
		Offset:  p * num,
		OrderBy: "sorted_id",
		Desc:    true,
	}); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	return
}
