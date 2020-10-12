package argus

import (
	"errors"
	"testing"
)

func TestLogHandler_StackTrace(t *testing.T) {
	logger := NewLogger()
	logger.StackTrace()
}

func TestLogHandler_ErrorLog(t *testing.T) {
	logger := NewLogger()
	var errs []Error
	errs = append(errs, Error{
		Err:    errors.New("test error"),
		Type:   BasicError,
		Values: map[string]interface{}{"test": 1},
	})
	logger.ErrorLog(errs)
}
