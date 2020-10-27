package private

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// ParseRegisterArticle parses request body to get article contents.
func ParseRegisterArticle(ctx *gin.Context, reqCh chan<- APIRegisterArticleReq, errCh chan<- error) {
	defer close(reqCh)

	// Parse request body.
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errCh <- err
		return
	}

	// Unmarshal body to json.
	var req APIRegisterArticleReq
	if err := json.Unmarshal(body, &req); err != nil {
		errCh <- err
		return
	}

	reqCh <- req
}

// ParseUpdateArticle parses request body to get article contents.
func ParseUpdateArticle(ctx *gin.Context, reqCh chan<- APIUpdateArticleReq, errCh chan<- error) {
	defer close(reqCh)

	// Parse request body.
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errCh <- err
		return
	}

	// Unmarshal body to json.
	var req APIUpdateArticleReq
	if err := json.Unmarshal(body, &req); err != nil {
		errCh <- err
		return
	}

	reqCh <- req
}

// ParseTogglePrivate parses request body.
func ParseTogglePrivate(ctx *gin.Context, reqCh chan<- APITogglePrivateReq, errCh chan<- error) {
	defer close(reqCh)

	// Parse request body.
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errCh <- err
		return
	}

	// Unmarshal body to json.
	var req APITogglePrivateReq
	if err := json.Unmarshal(body, &req); err != nil {
		errCh <- err
		return
	}

	reqCh <- req
}

// ParseRegisterDraft parses request body to get draft contents.
func ParseRegisterDraft(ctx *gin.Context, reqCh chan<- APIRegisterDraftReq, errCh chan<- error) {
	defer close(reqCh)

	// Parse request body.
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errCh <- err
		return
	}

	// Unmarshal body to json.
	var req APIRegisterDraftReq
	if err := json.Unmarshal(body, &req); err != nil {
		errCh <- err
		return
	}

	reqCh <- req
}

// ParseUpdateDraft parses request body to get draft contents.
func ParseUpdateDraft(ctx *gin.Context, reqCh chan<- APIUpdateDraftReq, errCh chan<- error) {
	defer close(reqCh)

	// Parse request body.
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errCh <- err
		return
	}

	// Unmarshal body to json.
	var req APIUpdateDraftReq
	if err := json.Unmarshal(body, &req); err != nil {
		errCh <- err
		return
	}

	reqCh <- req
}

// ParseDeleteDraft parses request body to get draft contents.
func ParseDeleteDraft(ctx *gin.Context, reqCh chan<- APIDeleteDraftReq, errCh chan<- error) {
	defer close(reqCh)

	// Parse request body.
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errCh <- err
		return
	}

	// Unmarshal body to json.
	var req APIDeleteDraftReq
	if err := json.Unmarshal(body, &req); err != nil {
		errCh <- err
		return
	}

	reqCh <- req
}

// ParseDeleteImage parses request body to get image contents.
func ParseDeleteImage(ctx *gin.Context, reqCh chan<- APIDeleteImageReq, errCh chan<- error) {
	defer close(reqCh)

	// Parse request body.
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		errCh <- err
		return
	}

	// Unmarshal body to json.
	var req APIDeleteImageReq
	if err := json.Unmarshal(body, &req); err != nil {
		errCh <- err
		return
	}

	reqCh <- req
}
