package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/repo"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDraftHandler(t *testing.T) {
	requestBody := `{
	"article": {
		"id": "TEST_ID",
		"title": "test",
		"categories": "c1&c2",
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
		assert.Equal(t, d.Title, "test")
		assert.Equal(t, d.Categories, "c1&c2")
		assert.Equal(t, d.ContentHash, "0123456789")
		assert.Equal(t, d.ImageHash, "9876543210")
		return
	}

	if err := DraftHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		*Errors = []argus.Error{}
		t.Fatalf("error happend in handler")
	}

	expectedBody := `{
	"id": "TEST_ID",
	"contentHash": "0123456789",
	"imageHash": "9876543210"
}`

	var buf bytes.Buffer

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)

	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Indent(&buf, body, "", "	"); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		*Errors = []argus.Error{}
		t.Fatalf("Unable to indent json string\n")
		return
	}
	assert.Equal(t, expectedBody, buf.String())
	buf.Reset()
}
