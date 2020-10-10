package private

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// ParseRegisterArticle parses request body to get article.
func ParseRegisterArticle(ctx *gin.Context, resc chan<- APIRegisterArticleReq, errc chan<- error) {
	defer close(resc)

	// Parse request body.
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errc <- err
		return
	}

	// Unmarshal body to json.
	var res APIRegisterArticleReq
	if err := json.Unmarshal(body, &res); err != nil {
		errc <- err
		return
	}

	resc <- res
}
