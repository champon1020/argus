package private

import (
	"net/http"

	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

// APIUpdateArticleReq is the request type.
type APIUpdateArticleReq struct {
	Article model.Article `json:"article"`
}

// APIUpdateArticle updates article on database.
func APIUpdateArticle(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for request.
	reqc := make(chan APIUpdateArticleReq, 1)

	// Channel for error variable.
	errc := make(chan error, 1)

	go ParseUpdateArticle(ctx, reqc, errc)

	req, ok := <-reqc
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errc
	}

	if err := db.UpdateArticle(&req.Article); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	ctx.AbortWithStatus(http.StatusOK)
	return nil
}
