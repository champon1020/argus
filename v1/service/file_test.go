package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/champon1020/argus/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetFileList(t *testing.T) {
	var (
		files []os.FileInfo
		err   error
	)

	dirPath := filepath.Join(argus.EnvVars.Get("resource"), "files")
	if files, err = GetFileList(dirPath); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happened")
	}

	assert.Equal(t, len(files), 3)
	assert.Equal(t, files[0].Name(), "file1.txt")
	assert.Equal(t, files[1].Name(), "file2.txt")
	assert.Equal(t, files[2].Name(), "file3.txt")
}

func TestOutputFile(t *testing.T) {
	dummyBody := "dummy"
	path := filepath.Join(argus.EnvVars.Get("resource"), "files", "a")
	defer func() {
		_ = DeleteFile(path)
	}()

	if err := OutputFile(path, []byte(dummyBody)); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happened")
	}

	dirPath := filepath.Join(argus.EnvVars.Get("resource"), "files")
	files, _ := GetFileList(dirPath)

	assert.Equal(t, len(files), 4)
	assert.Equal(t, files[0].Name(), "a")
}

func TestDeleteFile(t *testing.T) {
	path := filepath.Join(argus.EnvVars.Get("resource"), "files", "file1.txt")
	defer func() {
		_ = OutputFile(path, []byte(""))
	}()

	if err := DeleteFile(path); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happened")
	}

	dirPath := filepath.Join(argus.EnvVars.Get("resource"), "files")
	files, _ := GetFileList(dirPath)

	assert.Equal(t, len(files), 2)
}
