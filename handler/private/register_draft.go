package private

import (
	"net/http"

	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

// APIRegisterDraftReq is the request type.
type APIRegisterDraftReq struct {
	Draft model.Draft `json:"draft"`
}

// APIRegisterDraft is the private handler to register new draft to database.
func APIRegisterDraft(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for request.
	reqc := make(chan APIRegisterDraftReq)

	// Channel for error variable.
	errc := make(chan error)

	go ParseRegisterDraft(ctx, reqc, errc)

	req, ok := <-reqc
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errc
	}

	if err := db.RegisterDraft(&req.Draft); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	ctx.AbortWithStatus(http.StatusOK)
	return nil
}
