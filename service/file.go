package service

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
)

func ReadBody(r io.Reader) (body []byte, err error) {
	if body, err = ioutil.ReadAll(r); err != nil {
		IOReadError.SetErr(err).AppendTo(Errors)
	}
	return
}

// Update file.
// If not exist, create new file and save it.
func OutputFile(path string, body []byte) (err error) {
	var file *os.File
	if file, err = os.Create(path); err != nil {
		IOWriteError.SetErr(err).AppendTo(Errors)
		return
	}
	defer file.Close()
	file.Write(body)
	return
}

// Delete file.
// If not exist, add error.
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

// Save multipart form files.
func SaveMultipartFiles(path string, fileHeaders []*multipart.FileHeader) (err error) {
	var (
		f    multipart.File
		body []byte
	)
	for _, fh := range fileHeaders {
		if f, err = fh.Open(); err != nil {
			return
		}
		if body, err = ReadBody(f); err != nil {
			return
		}

		fn := path + fh.Filename
		if err = OutputFile(fn, body); err != nil {
			return
		}
	}
	return
}