package handler

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/champon1020/argus/v1"
	"github.com/champon1020/argus/v1/repo"
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
		"content": "TEST_CONTENT",
		"imageHash": "9876543210",
		"isPrivate": true
	}
}`

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/private/update/article/object",
		strings.NewReader(requestBody))

	repoCmdMock := func(_ repo.MySQL, a repo.Article) (_ error) {
		assert.Equal(t, a.Id, "TEST_ID")
		assert.Equal(t, a.Title, "test")
		assert.Equal(t, len(a.Categories), 1)
		assert.Equal(t, a.Categories[0].Id, "TEST_CA_ID")
		assert.Equal(t, a.Categories[0].Name, "c1")
		assert.Equal(t, a.Content, "TEST_CONTENT")
		assert.Equal(t, a.ImageHash, "9876543210")
		assert.Equal(t, a.Private, true)
		return
	}

	if err := UpdateArticleHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		*Errors = []argus.Error{}
		t.Fatalf("error happend in handler")
	}

	assert.Equal(t, w.Result().StatusCode, 200)
}
