package service_test

import (
	"os"
	"testing"
	"time"

	"github.com/champon1020/argus/service"
	"github.com/stretchr/testify/assert"
)

func TestResolveContentHash(t *testing.T) {
	contentHash := "article"
	fn := service.ResolveContentHash(contentHash)
	assert.Equal(t, "article", fn)

	contentHash = ""
	fn = service.ResolveContentHash(contentHash)
	today := time.Now()
	tStr := today.Format("20060102150405")
	assert.Equal(t, string([]rune(tStr)[:6]), string([]rune(fn)[:6]))
}

func TestResolveContentFilePath(t *testing.T) {
	contentHash := "test"
	dirName := "articles"
	path := service.ResolveContentFilePath(dirName, contentHash)
	assert.Equal(t, os.Getenv("ARGUS_RESOURCE_PATH")+"/articles/test", path)
}
