package private

import (
	"net/http"

	"github.com/champon1020/argus/handler"
	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

// APIFindArticles is the private handler to get all articles.
func APIFindArticles(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter p.
	pc := make(chan int, 1)

	// Channel for query parameter num.
	numc := make(chan int, 1)

	// Channel for error variable.
	errc := make(chan error, 2)

	// Response of this call.
	res := new(handler.APIFindArticlesRes)

	go handler.ParsePage(ctx, pc, errc)

	go handler.ParseNum(ctx, numc, errc)

	p, ok1 := <-pc
	num, ok2 := <-numc
	if !ok1 || !ok2 {
		return <-errc
	}

	doneFind := make(chan bool)
	doneCount := make(chan bool)

	go func() {
		defer close(doneFind)
		if err := db.FindAllArticles(
			&res.Articles,
			model.NewOp(num, (p-1)*num, "sorted_id", true),
		); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			errc <- err
			return
		}
		doneFind <- true
	}()

	go func() {
		defer close(doneCount)
		if err := db.CountAllArticles(
			&res.Count,
			model.NewOp(num, (p-1)*num, "sorted_id", true),
		); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			errc <- err
			return
		}
		doneCount <- true
	}()

	_, ok1 = <-doneFind
	_, ok2 = <-doneCount
	if !ok1 || !ok2 {
		return <-errc
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}

// APIFindArticleByIDRes is the response type.
type APIFindArticleByIDRes struct {
	Article model.Article `json:"article"`
}

// APIFindArticleByID is the private handler to get article by id.
func APIFindArticleByID(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter id.
	idc := make(chan string, 1)

	// Channel for error variable.
	errc := make(chan error, 1)

	// Response of this handler.
	res := new(APIFindArticleByIDRes)

	go handler.ParseID(ctx, idc, errc)

	id, ok := <-idc
	if !ok {
		return <-errc
	}

	if err := db.FindArticleByID(&res.Article, id); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}
