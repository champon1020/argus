package private

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// ParseRegisterArticle parses request body to get article contents.
func ParseRegisterArticle(ctx *gin.Context, reqc chan<- APIRegisterArticleReq, errc chan<- error) {
	defer close(reqc)

	// Parse request body.
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errc <- err
		return
	}

	// Unmarshal body to json.
	var req APIRegisterArticleReq
	if err := json.Unmarshal(body, &req); err != nil {
		errc <- err
		return
	}

	reqc <- req
}

// ParseUpdateArticle parses request body to get article contents.
func ParseUpdateArticle(ctx *gin.Context, reqc chan<- APIUpdateArticleReq, errc chan<- error) {
	defer close(reqc)

	// Parse request body.
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errc <- err
		return
	}

	// Unmarshal body to json.
	var req APIUpdateArticleReq
	if err := json.Unmarshal(body, &req); err != nil {
		errc <- err
		return
	}

	reqc <- req
}

// ParseRegisterDraft parses request body to get draft contents.
func ParseRegisterDraft(ctx *gin.Context, reqc chan<- APIRegisterDraftReq, errc chan<- error) {
	defer close(reqc)

	// Parse request body.
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errc <- err
		return
	}

	// Unmarshal body to json.
	var req APIRegisterDraftReq
	if err := json.Unmarshal(body, &req); err != nil {
		errc <- err
		return
	}

	reqc <- req
}
