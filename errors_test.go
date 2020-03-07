package argus

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultErr Error = Error{
	Err:  errors.New("test error"),
	Type: IOFailedReadError,
	Values: map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	},
}

func TestError_AppendTo(t *testing.T) {
	var errs []Error
	Error{
		Err:  errors.New("test error"),
		Type: IOFailedReadError,
	}.AppendTo(&errs)

	assert.Equal(t, 1, len(errs))
}

func TestError_String(t *testing.T) {
	actual := "test error"
	pre := defaultErr.String()

	assert.Equal(t, actual, pre)
}

func TestError_Marshal(t *testing.T) {
	if _, ok := defaultErr.Marshal(); ok != nil {
		t.Errorf("failed to marshal")
	}
}

func TestError_JSON(t *testing.T) {
	json := defaultErr.JSON()

	assert.Equal(t, map[string]interface{}{
		"Error": "test error",
		"Type":  IOFailedReadError,
		"Values": map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}, json)
}
