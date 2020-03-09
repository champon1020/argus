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

func TestDraftHandler(t *testing.T) {
	requestBody := `{
	"article": {
		"id": 1,
		"title": "test",
		"categories": [
			{"id": -1, "name": "c1"},
			{"id": -1, "name": "c2"}
		],
		"contentHash": "0123456789",
		"imageHash": "9876543210"
	},
	"contents": "<div>ok</div>"
}`

	defer func() {
		argus.StdLogger.ErrorLog(*Errors)
		service.DeleteFile(service.ResolveContentFilePath("0123456789", "drafts"))
	}()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"POST",
		"/draft/article",
		strings.NewReader(requestBody))

	repoCmdMock := func(_ repo.MySQL, d repo.Draft) (_ error) {
		assert.Equal(t, d.Id, 1)
		assert.Equal(t, d.Title, "test")
		assert.Equal(t, d.Categories, "c1&c2")
		assert.Equal(t, d.ContentHash, "0123456789")
		assert.Equal(t, d.ImageHash, "9876543210")
		return
	}

	var res *http.Response
	DraftHandler(ctx, repoCmdMock)
	res = w.Result()
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, len(*Errors), 0)
	*Errors = []argus.Error{}
}
