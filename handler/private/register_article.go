package private

import (
	"net/http"

	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

// APIRegisterArticleReq is the request type.
type APIRegisterArticleReq struct {
	Article model.Article `json:"article"`
	DraftID string        `json:"draftId"`
}

// APIRegisterArticle is the private hanlder to register new article to database.
func APIRegisterArticle(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for request.
	reqCh := make(chan APIRegisterArticleReq)

	// Channel for error variable.
	errCh := make(chan error)

	go ParseRegisterArticle(ctx, reqCh, errCh)

	req, ok := <-reqCh
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errCh
	}

	// Register new article to articles table.
	if err := db.RegisterArticle(&req.Article, req.DraftID); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	ctx.AbortWithStatus(http.StatusOK)
	return nil
}
