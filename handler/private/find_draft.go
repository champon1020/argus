package private

import (
	"net/http"

	"github.com/champon1020/argus/handler"
	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

// APIFindDraftByIDRes is the response type.
type APIFindDraftByIDRes struct {
	Draft model.Draft `json:"draft"`
}

// APIFindDraftByID is the private handler to get draft by id.
func APIFindDraftByID(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter id.
	idCh := make(chan string, 1)

	// Channel for error variable.
	errCh := make(chan error, 1)

	// Response of this handler.
	res := new(APIFindDraftByIDRes)

	go handler.ParseID(ctx, idCh, errCh)

	id, ok := <-idCh
	if !ok {
		return <-errCh
	}

	if err := db.FindDraftByID(&res.Draft, id); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	// Article is not exist.
	if res.Draft.ID == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return nil
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}

// APIFindDraftsRes is the response type.
type APIFindDraftsRes struct {
	Drafts []model.Draft `json:"drafts"`
	Count  int           `json:"count"`
}

// APIFindDrafts is the private handler to get all drafts.
func APIFindDrafts(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter p.
	pCh := make(chan int, 1)

	// Channel for query parameter num.
	numCh := make(chan int, 1)

	// Channel for error variable.
	errCh := make(chan error, 2)

	// Response of this handler.
	res := new(APIFindDraftsRes)

	go handler.ParsePage(ctx, pCh, errCh)

	go handler.ParseNum(ctx, numCh, errCh)

	p, ok1 := <-pCh
	num, ok2 := <-numCh
	if !ok1 || !ok2 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errCh
	}

	doneFind := make(chan bool)
	doneCount := make(chan bool)

	go func() {
		defer close(doneFind)
		if err := db.FindDrafts(
			&res.Drafts,
			model.NewOp(num, (p-1)*num, "id", true),
		); err != nil {
			errCh <- err
			return
		}
		doneFind <- true
	}()

	go func() {
		defer close(doneCount)
		if err := db.CountDrafts(&res.Count); err != nil {
			errCh <- err
			return
		}
		doneCount <- true
	}()

	_, ok1 = <-doneFind
	_, ok2 = <-doneCount
	if !ok1 || !ok2 {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return <-errCh
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}
