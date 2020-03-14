package handler

import (
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
		_ = service.DeleteFile(service.ResolveContentFilePath("0123456789", "drafts"))
	}()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"POST",
		"/api/draft/article",
		strings.NewReader(requestBody))

	repoCmdMock := func(_ repo.MySQL, d repo.Draft) (_ error) {
		assert.Equal(t, d.Id, 1)
		assert.Equal(t, d.Title, "test")
		assert.Equal(t, d.Categories, "c1&c2")
		assert.Equal(t, d.ContentHash, "0123456789")
		assert.Equal(t, d.ImageHash, "9876543210")
		return
	}

	if err := DraftHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)
}
