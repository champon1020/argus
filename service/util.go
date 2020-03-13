package service

import (
	"strings"
	"time"

	"github.com/champon1020/argus"
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
func ResolveContentFilePath(contentHash string, dirName string) string {
	fn := ResolveContentHash(contentHash)
	return argus.EnvVars.Get("resources") + "/" + dirName + "/" + fn
}

// Get only fine name from path(uri or url).
func ConvertPathToFileName(path string) string {
	if path == "" {
		return path
	}
	seps := strings.Split(path, "/")
	fileName := seps[len(seps)-1]
	return fileName
}
