package handler

import (
	"net/http"

	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

// APIFindCategoriesRes is the response type.
type APIFindCategoriesRes struct {
	Categories []model.Category `json:"categories"`
}

// APIFindCategories is the handler to get all categories with public articles.
func APIFindCategories(ctx *gin.Context, db model.DatabaseIface) error {
	res := new(APIFindCategoriesRes)

	if err := db.FindPublicCategories(&res.Categories, nil); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}
