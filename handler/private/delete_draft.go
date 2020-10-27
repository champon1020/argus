package private

import (
	"net/http"

	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

// APIDeleteDraftReq is the request type.
type APIDeleteDraftReq struct {
	ID string `json:"id"`
}

// APIDeleteDraft is the private handler to delete draft.
func APIDeleteDraft(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for request.
	reqCh := make(chan APIDeleteDraftReq, 1)

	// Channel for error variable.
	errCh := make(chan error, 1)

	go ParseDeleteDraft(ctx, reqCh, errCh)

	req, ok := <-reqCh
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errCh
	}

	if err := db.DeleteDraft(req.ID); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	ctx.AbortWithStatus(http.StatusOK)
	return nil
}
