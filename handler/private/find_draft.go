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
	idc := make(chan string, 1)

	// Channel for error variable.
	errc := make(chan error, 1)

	// Response of this handler.
	res := new(APIFindDraftByIDRes)

	go handler.ParseID(ctx, idc, errc)

	id, ok := <-idc
	if !ok {
		return <-errc
	}

	if err := db.FindDraftByID(&res.Draft, id); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
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
	pc := make(chan int, 1)

	// Channel for query parameter num.
	numc := make(chan int, 1)

	// Channel for error variable.
	errc := make(chan error, 2)

	// Response of this handler.
	res := new(APIFindDraftsRes)

	go handler.ParsePage(ctx, pc, errc)

	go handler.ParseNum(ctx, numc, errc)

	p, ok1 := <-pc
	num, ok2 := <-numc
	if !ok1 || !ok2 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errc
	}

	doneFind := make(chan bool)
	doneCount := make(chan bool)

	go func() {
		defer close(doneFind)
		if err := db.FindDrafts(
			&res.Drafts,
			model.NewOp(num, (p-1)*num, "sorted_id", true),
		); err != nil {
			errc <- err
			return
		}
		doneFind <- true
	}()

	go func() {
		defer close(doneCount)
		if err := db.CountDrafts(
			&res.Count,
			model.NewOp(num, (p-1)*num, "sorted_id", true),
		); err != nil {
			errc <- err
			return
		}
		doneCount <- true
	}()

	_, ok1 = <-doneFind
	_, ok2 = <-doneCount
	if !ok1 || !ok2 {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return <-errc
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}
