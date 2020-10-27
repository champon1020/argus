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

// APIUpdateArticle is the private handler to update article on database.
func APIUpdateArticle(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for request.
	reqCh := make(chan APIUpdateArticleReq, 1)

	// Channel for error variable.
	errCh := make(chan error, 1)

	go ParseUpdateArticle(ctx, reqCh, errCh)

	req, ok := <-reqCh
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errCh
	}

	if err := db.UpdateArticle(&req.Article); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	ctx.AbortWithStatus(http.StatusOK)
	return nil
}

// APITogglePrivateReq is the request type.
type APITogglePrivateReq struct {
	ID      string `json:"id"`
	Private bool   `json:"private"`
}

// APIUpdateIsPrivate is the private handler to update isPrivate property of article on database.
func APIUpdateIsPrivate(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for request
	reqCh := make(chan APITogglePrivateReq, 1)

	// Channel for error variable.
	errCh := make(chan error, 1)

	go ParseTogglePrivate(ctx, reqCh, errCh)

	req, ok := <-reqCh
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errCh
	}

	if err := db.UpdateArticlePrivate(req.ID, req.Private); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	ctx.AbortWithStatus(http.StatusOK)
	return nil
}
