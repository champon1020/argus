package handler

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDeleteImageController(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"GET",
		"/api/delete/image?imgName=image_test1.png",
		nil)

	defer func() {
		_ = service.OutputFile(
			filepath.Join(argus.EnvVars.Get("resource"), "images", "image_test1.png"),
			[]byte(""),
		)
	}()

	DeleteImageController(ctx)
	assert.Equal(t, w.Result().StatusCode, 200)
}

func TestDeleteImageHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"GET",
		"/api/delete/image?imgName=image_test1.png",
		nil)

	defer func() {
		_ = service.OutputFile(
			filepath.Join(argus.EnvVars.Get("resource"), "images", "image_test1.png"),
			[]byte(""),
		)
	}()

	if err := DeleteImageHandler(ctx); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}
}
