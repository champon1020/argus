package handler

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/champon1020/argus/v1/repo"

	"github.com/champon1020/argus/v1"
	"github.com/champon1020/argus/v1/service"
	"github.com/gin-gonic/gin"
)

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

func TestDeleteDraftHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"GET",
		"/api/private/delete/draft?id=draft_test1&contentHash=draft_test1",
		nil)

	defer func() {
		_ = service.OutputFile(
			filepath.Join(argus.EnvVars.Get("resource"), "drafts", "draft_test1_md"),
			[]byte(""),
		)
	}()

	repoCmd := func(_ repo.MySQL, _ repo.Draft) error {
		return nil
	}

	if err := DeleteDraftHandler(ctx, repoCmd); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}
}
