package private

import (
	"net/http"

	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

// APIUpdateDraftReq is the request type.
type APIUpdateDraftReq struct {
	Draft model.Draft `json:"draft"`
}

// APIUpdateDraft updates draft on database.
func APIUpdateDraft(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for request.
	reqc := make(chan APIUpdateDraftReq, 1)

	// Channel for error variable.
	errc := make(chan error, 1)

	go ParseUpdateDraft(ctx, reqc, errc)

	req, ok := <-reqc
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errc
	}

	if err := db.UpdateDraft(&req.Draft); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	ctx.AbortWithStatus(http.StatusOK)
	return nil
}
