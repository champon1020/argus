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

// APIRegisterDraftRes is the response type.
type APIRegisterDraftRes struct {
	ID        string `json:"id"`
	content   string `json:"content"`
	imageHash string `json:"imageHash"`
}

// APIRegisterDraft is the private handler to register new draft to database.
func APIRegisterDraft(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for request.
	reqCh := make(chan APIRegisterDraftReq)

	// Channel for error variable.
	errCh := make(chan error)

	go ParseRegisterDraft(ctx, reqCh, errCh)

	req, ok := <-reqCh
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errCh
	}

	// Generate new draft id.
	req.Draft.ID = model.GetNewID(model.TypeDraft)

	if err := db.RegisterDraft(&req.Draft); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	res := APIRegisterDraftRes{ID: req.Draft.ID}

	ctx.JSON(http.StatusOK, res)
	return nil
}
