package handler

import (
	"net/http"

	"github.com/champon1020/argus/v2/model"
	"github.com/gin-gonic/gin"
)

// ArticlesResponse is the response type.
type ArticlesResponse struct {
	Articles []model.Article `json:"articles"`
	Count    int             `json:"count"`
}

// FindArticlesList gets all public articles.
func FindArticlesList(ctx *gin.Context, db model.DatabaseIface) error {
	// Channel for query parameter p.
	pc := make(chan int, 1)

	// Channel for query parameter num.
	numc := make(chan int, 1)

	// Channel for error variable.
	errc := make(chan error, 3)

	// Response of this call.
	res := new(ArticlesResponse)

	// Parse page number from gin context.
	go func() {
		defer close(pc)
		if p, err := ParsePage(ctx); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			errc <- err
		} else {
			pc <- p
		}
	}()

	// Parse the number of articles to response to client from gin context.
	go func() {
		defer close(numc)
		if num, err := ParseNum(ctx); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			errc <- err
		} else {
			numc <- num
		}
	}()

	donec := make(chan bool)

	// Search for articles.
	go func() {
		p, ok1 := <-pc
		num, ok2 := <-numc

		if !ok1 || !ok2 {
			return
		}

		if err := db.FindPublicArticles(&res.Articles, &model.QueryOptions{
			Limit:   num,
			Offset:  (p - 1) * num,
			OrderBy: "sorted_id",
			Desc:    true,
		}); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			errc <- err
			return
		}
		donec <- true
	}()

	// Count the nubmer of articles.
	go func() {
		if err := db.CountArticles(&res.Count); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			errc <- err
			return
		}
		donec <- true
	}()

	var done int

	// If some errors are occurred, channel errc receives errors
	// and return error variable.
	// Or, if database calls are done certain times, return the response.
	for {
		select {
		case err := <-errc:
			return err
		case <-donec:
			done++
			if done == 2 {
				// Set response with status 200.
				ctx.JSON(http.StatusOK, res)
				return nil
			}
		}
	}
}
