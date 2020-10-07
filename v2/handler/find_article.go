package handler

import (
	"net/http"

	"github.com/champon1020/argus/v2/model"
	"github.com/gin-gonic/gin"
)

// APIFindArticlesRes is the response type.
type APIFindArticlesRes struct {
	Articles []model.Article `json:"articles"`
	Count    int             `json:"count"`
}

// FindArticles gets all public articles.
func FindArticles(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter p.
	pc := make(chan int, 1)

	// Channel for query parameter num.
	numc := make(chan int, 1)

	// Channel for error variable.
	errc := make(chan error, 2)

	// Response of this call.
	res := new(APIFindArticlesRes)

	// Parse page number from gin context.
	go ParsePage(ctx, pc, errc)

	// Parse the number of articles to response to client from gin context.
	go ParseNum(ctx, numc, errc)

	p, ok1 := <-pc
	num, ok2 := <-numc
	if !ok1 || !ok2 {
		err := <-errc
		return err
	}

	doneCount := make(chan bool)
	doneFind := make(chan bool)

	// Search for articles.
	go func() {
		if err := db.FindPublicArticles(
			&res.Articles,
			model.NewOp(num, (p-1)*num, "sorted_id", true),
		); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			errc <- err
			return
		}
		doneFind <- true
	}()

	// Count the nubmer of articles.
	go func() {
		if err := db.CountPublicArticles(
			&res.Count,
			model.NewOp(num, (p-1)*num, "sorted_id", true),
		); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			errc <- err
			return
		}
		doneCount <- true
	}()

	// If some errors are occurred, channel errc receives errors
	// and return error variable.
	// Or, if all database calls are done, return the response.
	_, ok1 = <-doneFind
	_, ok2 = <-doneCount
	if !ok1 || !ok2 {
		err := <-errc
		return err
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}

// FindArticlesByTitle gets public articles
// whose title is specified at query parameter.
func FindArticlesByTitle(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter p.
	pc := make(chan int, 1)

	// Channel for query parameter num.
	numc := make(chan int, 1)

	// Channel for query parameter title.
	titlec := make(chan string, 1)

	// Channel for error variable.
	errc := make(chan error, 3)

	// Response of this call.
	res := new(APIFindArticlesRes)

	// Parse page number from gin context.
	go ParsePage(ctx, pc, errc)

	// Parse the number of articles to response to client from gin context.
	go ParseNum(ctx, numc, errc)

	// Parse article title from gin context.
	go ParseTitle(ctx, titlec, errc)

	p, ok1 := <-pc
	num, ok2 := <-numc
	title, ok3 := <-titlec
	if !ok1 || !ok2 || !ok3 {
		err := <-errc
		return err
	}

	doneCount := make(chan bool)
	doneFind := make(chan bool)

	// Search for articles by title.
	go func() {
		if err := db.FindPublicArticlesByTitle(
			&res.Articles,
			title,
			model.NewOp(p, (p-1)*num, "sorted_id", true),
		); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			errc <- err
			return
		}
		doneFind <- true
	}()

	// Count the number of articles with title.
	go func() {
		if err := db.CountPublicArticlesByTitle(
			&res.Count,
			title,
			model.NewOp(p, (p-1)*num, "sorted_id", true),
		); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			errc <- err
			return
		}
		doneCount <- true
	}()

	// If some errors are occurred, channel errc receives errors
	// and return error variable.
	// Or, if database calls are done certain times, return the response.
	_, ok1 = <-doneFind
	_, ok2 = <-doneCount
	if !ok1 || !ok2 {
		err := <-errc
		return err
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}
