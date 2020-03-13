package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/champon1020/argus"
	repo "github.com/champon1020/argus/repository"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterArticleHandler(t *testing.T) {
	requestBody := `{
	"article": {
		"id": 1,
		"title": "test",
		"categories": [
			{"id": 1, "name": "c1"}
		],
		"contentHash": "0123456789",
		"imageHash": "9876543210",
		"private": false
	},
	"contents": "<div>ok</div>"
}`

	defer func() {
		argus.StdLogger.ErrorLog(*Errors)
		service.DeleteFile(service.ResolveContentFilePath("0123456789", "articles"))
	}()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"POST",
		"/register/article",
		strings.NewReader(requestBody))

	repoCmdMock := func(_ repo.MySQL, a repo.Article) (_ error) {
		assert.Equal(t, a.Id, 1)
		assert.Equal(t, a.Title, "test")
		assert.Equal(t, len(a.Categories), 1)
		assert.Equal(t, a.Categories[0].Name, "c1")
		assert.Equal(t, a.ContentHash, "0123456789")
		assert.Equal(t, a.ImageHash, "9876543210")
		assert.Equal(t, a.Private, false)
		return
	}

	var res *http.Response
	RegisterArticleHandler(ctx, repoCmdMock)
	res = w.Result()
	assert.Equal(t, res.StatusCode, 200)

	// see details in debug.log
	if len(*Errors) != 0 {
		argus.Logger.ErrorLog(*Errors)
	}
	assert.Equal(t, len(*Errors), 0)
	*Errors = []argus.Error{}
}
