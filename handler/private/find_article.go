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
	pCh := make(chan int, 1)

	// Channel for query parameter num.
	numCh := make(chan int, 1)

	// Channel for error variable.
	errCh := make(chan error, 2)

	// Response of this call.
	res := new(handler.APIFindArticlesRes)

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
		if err := db.FindAllArticles(
			&res.Articles,
			model.NewOp(num, (p-1)*num, "id", true),
		); err != nil {
			errCh <- err
			return
		}
		doneFind <- true
	}()

	go func() {
		defer close(doneCount)
		if err := db.CountAllArticles(&res.Count); err != nil {
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

// APIFindArticleByIDRes is the response type.
type APIFindArticleByIDRes struct {
	Article model.Article `json:"article"`
}

// APIFindArticleByID is the private handler to get article by id.
func APIFindArticleByID(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter id.
	idCh := make(chan string, 1)

	// Channel for error variable.
	errCh := make(chan error, 1)

	// Response of this handler.
	res := new(APIFindArticleByIDRes)

	go handler.ParseID(ctx, idCh, errCh)

	id, ok := <-idCh
	if !ok {
		return <-errCh
	}

	if err := db.FindArticleByID(&res.Article, id); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return err
	}

	// Article is not exist.
	if res.Article.ID == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return nil
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}
