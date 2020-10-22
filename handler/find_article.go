package handler

import (
	"net/http"
	"sync"

	"github.com/champon1020/argus/model"
	"github.com/gin-gonic/gin"
)

// APIFindArticlesRes is the response type.
type APIFindArticlesRes struct {
	Articles []model.Article `json:"articles"`
	Count    int             `json:"count"`
}

// APIFindArticles is the handler to get all public articles.
func APIFindArticles(ctx *gin.Context, db model.DatabaseIface) error {
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
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errc
	}

	doneFind := make(chan bool)
	doneCount := make(chan bool)

	// Search for articles.
	go func() {
		defer close(doneFind)
		if err := db.FindPublicArticles(
			&res.Articles,
			model.NewOp(num, (p-1)*num, "sorted_id", true),
		); err != nil {
			errc <- err
			return
		}
		doneFind <- true
	}()

	// Count the nubmer of articles.
	go func() {
		defer close(doneCount)
		if err := db.CountPublicArticles(&res.Count); err != nil {
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
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return <-errc
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}

// APIFindArticlesBySortedIDRes is the response type.
type APIFindArticlesBySortedIDRes struct {
	Article     model.Article `json:"article"`
	PrevArticle model.Article `json:"prevArticle"`
	NextArticle model.Article `json:"nextArticle"`
}

// APIFindArticlesBySortedID is the handler
// to get public articles whose sorted id
// is greater than and equal to thet of query parameter.
func APIFindArticlesBySortedID(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter sortedID.
	sidc := make(chan int, 1)

	// Channel for error variable.
	errc := make(chan error, 3)

	// Response of this call.
	res := new(APIFindArticlesBySortedIDRes)

	// Parse article sorted id from gin context.
	go ParseSortedID(ctx, sidc, errc)

	sortedID, ok := <-sidc
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errc
	}

	ac := make(chan []model.Article)

	// Search for the articles by sorted id.
	go func() {
		defer close(ac)
		var a []model.Article
		if err := db.FindPublicArticlesGeSortedID(
			&a,
			sortedID-1,
			model.NewOp(3, 0, "sorted_id", true),
		); err != nil {
			errc <- err
			return
		}
		ac <- a
	}()

	articles, ok := <-ac
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return <-errc
	}

	// Assing fetched articles to response type.
	// If the sorted id is equal to that of query parameter,
	// it's assigned to Article Field.
	//
	// If the sorted id is less than that of query parameter,
	// its title assigned to prevTitle Field.
	//
	// If the sorted id is greater than that of query parameter,
	// its title assigned to nextTitle Field.
	wg := new(sync.WaitGroup)
	for _, a := range articles {
		wg.Add(1)
		go func(a model.Article) {
			defer wg.Done()
			if a.SortedID == sortedID {
				res.Article = a
			}
			if a.SortedID < sortedID {
				res.PrevArticle = a
			}
			if a.SortedID > sortedID {
				res.NextArticle = a
			}
		}(a)
	}
	wg.Wait()

	ctx.JSON(http.StatusOK, res)

	return nil
}

// APIFindArticlesByTitle is the handler to get public articles by title.
func APIFindArticlesByTitle(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter p.
	pc := make(chan int, 1)

	// Channel for query parameter num.
	numc := make(chan int, 1)

	// Channel for query parameter title.
	titlec := make(chan string, 1)

	// Channel for error variable.
	errc := make(chan error, 2)

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
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errc
	}

	doneFind := make(chan bool)
	doneCount := make(chan bool)

	// Search for articles by title.
	go func() {
		defer close(doneFind)
		if err := db.FindPublicArticlesByTitle(
			&res.Articles,
			title,
			model.NewOp(num, (p-1)*num, "sorted_id", true),
		); err != nil {
			errc <- err
			return
		}
		doneFind <- true
	}()

	// Count the number of articles with title.
	go func() {
		defer close(doneCount)
		if err := db.CountPublicArticlesByTitle(
			&res.Count,
			title,
		); err != nil {
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
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return <-errc
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}

// APIFindArticlesByCategory is the handler to get puclic articles by category id.
func APIFindArticlesByCategory(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter p.
	pc := make(chan int, 1)

	// Channel for query parameter num.
	numc := make(chan int, 1)

	// Channel for query parameter categoryID.
	catec := make(chan string, 1)

	// Channel for error variable.
	errc := make(chan error, 2)

	// Response of this call.
	res := new(APIFindArticlesRes)

	go ParsePage(ctx, pc, errc)

	go ParseNum(ctx, numc, errc)

	go ParseCategoryID(ctx, catec, errc)

	p, ok1 := <-pc
	num, ok2 := <-numc
	categoryID, ok3 := <-catec
	if !ok1 || !ok2 || !ok3 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errc
	}

	doneFind := make(chan bool)
	doneCount := make(chan bool)

	go func() {
		defer close(doneFind)
		if err := db.FindPublicArticlesByCategory(
			&res.Articles,
			categoryID,
			model.NewOp(num, (p-1)*num, "sorted_id", true),
		); err != nil {
			errc <- err
			return
		}
		doneFind <- true
	}()

	go func() {
		defer close(doneCount)
		if err := db.CountPublicArticlesByCategory(
			&res.Count,
			categoryID,
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
