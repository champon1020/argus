package handler

import (
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/repo"
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
		_ = service.DeleteFile(service.ResolveContentFilePath("0123456789", "articles"))
	}()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"POST",
		"/api/register/article",
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

	if err := RegisterArticleHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)
}

func TestRegisterImageHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// prepare dummy image
	//mFileHeader := multipart.FileHeader{
	//	Filename: "test_image.png",
	//	Header:   make(map[string][]string),
	//	Size:     16,
	//}

	// create test request
	ctx.Request = httptest.NewRequest(
		"POST",
		"/register/article",
		nil)
	ctx.Request.Header.Set("Content-Type", "multipart/form-data")
	ctx.Request.MultipartForm = &multipart.Form{
		//Value: make(map[string][]string),
		//File: map[string][]*multipart.FileHeader{
		//	"images": {&mFileHeader},
		//},
	}

	if err := RegisterImageHandler(ctx); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)
}
