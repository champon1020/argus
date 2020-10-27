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
	pCh := make(chan int, 1)

	// Channel for query parameter num.
	numCh := make(chan int, 1)

	// Channel for error variable.
	errCh := make(chan error, 2)

	// Response of this call.
	res := new(APIFindArticlesRes)

	// Parse page number from gin context.
	go ParsePage(ctx, pCh, errCh)

	// Parse the number of articles to response to client from gin context.
	go ParseNum(ctx, numCh, errCh)

	p, ok1 := <-pCh
	num, ok2 := <-numCh
	if !ok1 || !ok2 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errCh
	}

	doneFind := make(chan bool)
	doneCount := make(chan bool)

	// Search for articles.
	go func() {
		defer close(doneFind)
		if err := db.FindPublicArticles(
			&res.Articles,
			model.NewOp(num, (p-1)*num, "id", true),
		); err != nil {
			errCh <- err
			return
		}
		doneFind <- true
	}()

	// Count the nubmer of articles.
	go func() {
		defer close(doneCount)
		if err := db.CountPublicArticles(&res.Count); err != nil {
			errCh <- err
			return
		}
		doneCount <- true
	}()

	// If some errors are occurred, channel errCh receives errors
	// and return error variable.
	// Or, if all database calls are done, return the response.
	_, ok1 = <-doneFind
	_, ok2 = <-doneCount
	if !ok1 || !ok2 {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return <-errCh
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}

// APIFindArticlesByIDRes is the response type.
type APIFindArticlesByIDRes struct {
	Article     model.Article `json:"article"`
	PrevArticle model.Article `json:"prevArticle"`
	NextArticle model.Article `json:"nextArticle"`
}

// APIFindArticlesByID is the handler to get public articles
// whose id is greater than and equal to thet of query parameter.
func APIFindArticlesByID(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter sortedID.
	idCh := make(chan string, 1)

	// Channel for error variable.
	errCh := make(chan error, 3)

	// Response of this call.
	res := new(APIFindArticlesByIDRes)

	// Parse article id from gin context.
	go ParseID(ctx, idCh, errCh)

	id, ok := <-idCh
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errCh
	}

	aCh1 := make(chan []model.Article, 1)
	aCh2 := make(chan []model.Article, 1)

	// Search for an article whose id is equal to query parameter's.
	// Also searchs for an article whose id is greater than query parameter's.
	go func() {
		defer close(aCh1)
		var a []model.Article
		if err := db.FindPublicArticleGeID(
			&a,
			id,
			model.NewOp(2, 0, "", false),
		); err != nil {
			errCh <- err
			return
		}
		aCh1 <- a
	}()

	// Searchs for an article whose id is less than query parameter's.
	go func() {
		defer close(aCh2)
		var a []model.Article
		if err := db.FindPublicArticleLeID(
			&a,
			id,
			model.NewOp(1, 0, "", false),
		); err != nil {
			errCh <- err
			return
		}
		aCh2 <- a
	}()

	articles, ok1 := <-aCh1
	articles2, ok2 := <-aCh2
	if !ok1 || !ok2 {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return <-errCh
	}

	articles = append(articles, articles2...)

	// Assing fetched articles to response type.
	// If the id is equal to that of query parameter,
	// it's assigned to Article Field.
	//
	// If the id is less than that of query parameter,
	// its title assigned to prevTitle Field.
	//
	// If the id is greater than that of query parameter,
	// its title assigned to nextTitle Field.
	wg := new(sync.WaitGroup)
	for _, a := range articles {
		wg.Add(1)
		go func(a model.Article) {
			defer wg.Done()
			if a.ID == id {
				res.Article = a
			}
			if a.ID < id {
				res.PrevArticle = a
			}
			if a.ID > id {
				res.NextArticle = a
			}
		}(a)
	}
	wg.Wait()

	// Article is not exist.
	if res.Article.ID == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return nil
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}

// APIFindArticlesByTitle is the handler to get public articles by title.
func APIFindArticlesByTitle(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter p.
	pCh := make(chan int, 1)

	// Channel for query parameter num.
	numCh := make(chan int, 1)

	// Channel for query parameter title.
	titleCh := make(chan string, 1)

	// Channel for error variable.
	errCh := make(chan error, 2)

	// Response of this call.
	res := new(APIFindArticlesRes)

	// Parse page number from gin context.
	go ParsePage(ctx, pCh, errCh)

	// Parse the number of articles to response to client from gin context.
	go ParseNum(ctx, numCh, errCh)

	// Parse article title from gin context.
	go ParseTitle(ctx, titleCh, errCh)

	p, ok1 := <-pCh
	num, ok2 := <-numCh
	title, ok3 := <-titleCh
	if !ok1 || !ok2 || !ok3 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errCh
	}

	doneFind := make(chan bool)
	doneCount := make(chan bool)

	// Search for articles by title.
	go func() {
		defer close(doneFind)
		if err := db.FindPublicArticlesByTitle(
			&res.Articles,
			title,
			model.NewOp(num, (p-1)*num, "id", true),
		); err != nil {
			errCh <- err
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
			errCh <- err
			return
		}
		doneCount <- true
	}()

	// If some errors are occurred, channel errCh receives errors
	// and return error variable.
	// Or, if database calls are done certain times, return the response.
	_, ok1 = <-doneFind
	_, ok2 = <-doneCount
	if !ok1 || !ok2 {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return <-errCh
	}

	ctx.JSON(http.StatusOK, res)

	return nil
}

// APIFindArticlesByCategory is the handler to get puclic articles by category id.
func APIFindArticlesByCategory(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter p.
	pCh := make(chan int, 1)

	// Channel for query parameter num.
	numCh := make(chan int, 1)

	// Channel for query parameter categoryID.
	cateCh := make(chan string, 1)

	// Channel for error variable.
	errCh := make(chan error, 2)

	// Response of this call.
	res := new(APIFindArticlesRes)

	go ParsePage(ctx, pCh, errCh)

	go ParseNum(ctx, numCh, errCh)

	go ParseCategoryID(ctx, cateCh, errCh)

	p, ok1 := <-pCh
	num, ok2 := <-numCh
	categoryID, ok3 := <-cateCh
	if !ok1 || !ok2 || !ok3 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return <-errCh
	}

	doneFind := make(chan bool)
	doneCount := make(chan bool)

	go func() {
		defer close(doneFind)
		if err := db.FindPublicArticlesByCategory(
			&res.Articles,
			categoryID,
			model.NewOp(num, (p-1)*num, "id", true),
		); err != nil {
			errCh <- err
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
