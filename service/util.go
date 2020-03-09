package service

import (
	"os"
	"strings"
	"time"

	"github.com/champon1020/argus"
)

var (
	Logger        = argus.Logger
	Errors        = &argus.Errors
	IOWriteError  = argus.NewError(argus.IOFailedWriteError)
	IORemoveError = argus.NewError(argus.IOFailedRemoveError)
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
	basePath := os.Getenv("ARGUS_ARTICLE")
	return basePath + "/" + dirName + "/" + fn
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

func OutputFile(path string, content string) (err error) {
	var file *os.File
	if file, err = os.Create(path); err != nil {
		IOWriteError.SetErr(err).AppendTo(Errors)
		return
	}
	defer file.Close()
	file.Write(([]byte)(content))
	return
}

func DeleteFile(path string) (err error) {
	if _, err := os.Stat(path); err != nil {
		Logger.Println("No such file: [handler.util] DeleteFile()")
		return err
	}
	if err := os.Remove(path); err != nil {
		IORemoveError.SetErr(err).AppendTo(Errors)
		return err
	}
	return
}
