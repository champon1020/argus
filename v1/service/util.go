package service

import (
	"path/filepath"
	"time"

	"github.com/champon1020/argus/v1"
)

var (
	Errors               = &argus.Errors
	IOOpenError          = argus.NewError(argus.IOFailedOpenError)
	IOReadError          = argus.NewError(argus.IOFailedReadError)
	IOWriteError         = argus.NewError(argus.IOFailedWriteError)
	IORemoveError        = argus.NewError(argus.IOFailedRemoveError)
	IOCloseError         = argus.NewError(argus.IOFailedCloseError)
	MultiFormatOpenError = argus.NewError(argus.MultiFormatFailedOpenError)
)

// If contentHash is empty, generate hash from date.
// If not, return.
func ResolveContentHash(contentHash string) string {
	if contentHash == "" {
		t := time.Now()
		return t.Format("20060102150405")
	}
	return contentHash
}

// Get file path uri from hash(fine name) and dir name.
func ResolveContentFilePath(dirName string, contentHash string) string {
	fn := ResolveContentHash(contentHash)
	return filepath.Join(argus.EnvVars.Get("resource"), dirName, fn)
}
