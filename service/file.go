package service

import (
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

func ReadBody(r io.Reader) (body []byte, err error) {
	if body, err = ioutil.ReadAll(r); err != nil {
		IOReadError.SetErr(err).AppendTo(Errors)
	}
	return
}

func GetFileList(dirPath string) (files []os.FileInfo, err error) {
	if files, err = ioutil.ReadDir(dirPath); err != nil {
		IOReadError.SetErr(err).AppendTo(Errors)
	}
	return
}

// Update file.
// If not exist, create new file and save it.
func OutputFile(path string, body []byte) (err error) {
	var file *os.File
	if file, err = os.Create(path); err != nil {
		IOOpenError.SetErr(err).AppendTo(Errors)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			IOCloseError.SetErr(err).AppendTo(Errors)
		}
	}()
	if _, err = file.Write(body); err != nil {
		IOWriteError.SetErr(err).AppendTo(Errors)
	}
	return
}

// Delete file.
// If not exist, add error.
func DeleteFile(path string) (err error) {
	if _, err := os.Stat(path); err != nil {
		IORemoveError.
			SetErr(errors.New("no such file: DeleteFile()")).
			SetValues("path", path).
			AppendTo(Errors)
		return err
	}
	if err := os.Remove(path); err != nil {
		IORemoveError.SetErr(err).AppendTo(Errors)
		return err
	}
	return
}

// Save multipart form files.
func SaveMultipartFiles(path string, fileHeaders []*multipart.FileHeader) (err error) {
	var (
		f    multipart.File
		body []byte
	)
	for _, fh := range fileHeaders {
		if f, err = fh.Open(); err != nil {
			MultiFormatOpenError.
				SetValues("file", f).
				SetErr(err).
				AppendTo(Errors)
			return
		}
		if body, err = ReadBody(f); err != nil {
			return
		}

		fn := filepath.Join(path, fh.Filename)
		if err = OutputFile(fn, body); err != nil {
			return
		}
	}
	return
}
