package service

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResolveContentHash(t *testing.T) {
	contentHash := "article"
	fn := ResolveContentHash(contentHash)
	assert.Equal(t, "article", fn)

	contentHash = ""
	fn = ResolveContentHash(contentHash)
	today := time.Now()
	tStr := today.Format("20060102150405")
	assert.Equal(t, string([]rune(tStr)[:6]), string([]rune(fn)[:6]))
}

func TestResolveContentFilePath(t *testing.T) {
	contentHash := "test"
	dirName := "articles"
	path := ResolveContentFilePath(contentHash, dirName)
	assert.Equal(t, os.Getenv("ARGUS_ARTICLE")+"/articles/test", path)
}

func TestConvertPathToFileName(t *testing.T) {
	path := "/this/is/test/article"
	fn := ConvertPathToFileName(path)
	assert.Equal(t, "article", fn)
}
