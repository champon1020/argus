package argus

import (
	"encoding/json"
	"errors"
)

// Error is the base error type of this project
// which implements error interface.
type Error struct {
	// Msg is custom error.
	Msg error

	// Err is default error.
	Err error

	// Values is for helping debug.
	Values map[string]interface{}
}

func (e Error) Error() string {
	return e.Err.Error()
}

// AppendValue adds new value for helping debug.
func (e *Error) AppendValue(key string, val interface{}) *Error {
	e.Values[key] = val
	return e
}

// MarshalIndent makes json indentation format.
func (e *Error) MarshalIndent() ([]byte, error) {
	return json.MarshalIndent(e.JSON(), "", "  ")
}

// JSON parses Error to json format.
func (e *Error) JSON() interface{} {
	if e.Err == nil {
		e.Err = errors.New("nil")
	}

	jsonMap := map[string]interface{}{
		"Msg":    e.Msg.Error(),
		"Error":  e.Err.Error(),
		"Values": e.Values,
	}

	return jsonMap
}

// Log outputs error log as json.
func (e *Error) Log() {
	bytes, _ := e.MarshalIndent()
	Logger.Printf("%s\n", string(bytes))
}

// NewError creates new error.
func NewError(msg error, err error) *Error {
	return &Error{
		Msg:    msg,
		Err:    err,
		Values: make(map[string]interface{}),
	}
}
