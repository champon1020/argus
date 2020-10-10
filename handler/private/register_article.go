package private

import (
	"net/http"

	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

// APIRegisterArticleReq is the request type.
type APIRegisterArticleReq struct {
	Article model.Article `json:"article"`
}

// APIRegisterArticle registers new article to database.
func APIRegisterArticle(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for article object.
	resc := make(chan APIRegisterArticleReq)

	// Channel for error variable.
	errc := make(chan error)

	go ParseRegisterArticle(ctx, resc, errc)

	res, ok := <-resc
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errc
	}

	if err := db.RegisterArticle(&res.Article); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	ctx.AbortWithStatus(http.StatusOK)
	return nil
}
