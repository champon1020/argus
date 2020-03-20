package handler

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/repo"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateArticleHandler(t *testing.T) {
	requestBody := `{
	"article": {
		"id": "TEST_ID",
		"title": "test",
		"categories": [
			{"id": "TEST_CA_ID", "name": "c1"}
		],
		"contentHash": "0123456789",
		"imageHash": "9876543210",
		"private": false
	},
	"contents": "<div>ok</div>"
}`

	defer func() {
		_ = service.DeleteFile(service.ResolveContentFilePath("0123456789", "articles"))
	}()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/update/article",
		strings.NewReader(requestBody))

	repoCmdMock := func(_ repo.MySQL, a repo.Article) (_ error) {
		assert.Equal(t, a.Id, "TEST_ID")
		assert.Equal(t, a.Title, "test")
		assert.Equal(t, len(a.Categories), 1)
		assert.Equal(t, a.Categories[0].Id, "TEST_CA_ID")
		assert.Equal(t, a.Categories[0].Name, "c1")
		assert.Equal(t, a.ContentHash, "0123456789")
		assert.Equal(t, a.ImageHash, "9876543210")
		assert.Equal(t, a.Private, false)
		return
	}

	if err := UpdateArticleHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		*Errors = []argus.Error{}
		t.Fatalf("error happend in handler")
	}

	assert.Equal(t, w.Result().StatusCode, 200)
}

func TestUpdateArticleObjController(t *testing.T) {
	requestBody := `{
	"article": {
		"id": "TEST_ID",
		"title": "test",
		"categories": [
			{"id": "TEST_CA_ID", "name": "c1"}
		],
		"contentHash": "0123456789",
		"imageHash": "9876543210",
		"isPrivate": true
	}
}`

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/update/article",
		strings.NewReader(requestBody))

	repoCmdMock := func(_ repo.MySQL, a repo.Article) (_ error) {
		assert.Equal(t, a.Id, "TEST_ID")
		assert.Equal(t, a.Title, "test")
		assert.Equal(t, len(a.Categories), 1)
		assert.Equal(t, a.Categories[0].Id, "TEST_CA_ID")
		assert.Equal(t, a.Categories[0].Name, "c1")
		assert.Equal(t, a.ContentHash, "0123456789")
		assert.Equal(t, a.ImageHash, "9876543210")
		assert.Equal(t, a.Private, true)
		return
	}

	if err := UpdateArticleObjHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		*Errors = []argus.Error{}
		t.Fatalf("error happend in handler")
	}

	assert.Equal(t, w.Result().StatusCode, 200)
}
