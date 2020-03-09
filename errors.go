package argus

import (
	"encoding/json"
)

type ErrorType string

const (
	// database error types
	DbRuntimeError     ErrorType = "DbRuntimeError"
	DbScanFailedError  ErrorType = "DbScanFailedError"
	DbCmdFailedError   ErrorType = "DbCmdFailedError"
	DbQueryFailedError ErrorType = "DbQueryFailedError"
	DbCloseFailedError ErrorType = "DbCloseFailedError"

	// io or json error types
	IOFailedReadError      ErrorType = "IOFailedReadError"
	IOFailedWriteError     ErrorType = "IOFailedWriteError"
	IOFailedRemoveError    ErrorType = "IOFailedRemoveError"
	IOFailedMarshalError   ErrorType = "IOFailedMarshalError"
	IOFailedUnmarshalError ErrorType = "IOFailedUnmarshalError"
)

type Error struct {
	Err    error
	Type   ErrorType
	Values map[string]interface{}
}

var Errors []Error

func NewError(t ErrorType) Error {
	return Error{Type: t, Values: make(map[string]interface{})}
}

func (e Error) AppendTo(errs *[]Error) {
	*errs = append(*errs, e)
}

func (e *Error) SetErr(err error) *Error {
	e.Err = err
	return e
}

func (e *Error) SetValues(key string, value interface{}) *Error {
	e.Values[key] = value
	return e
}

func (e Error) String() string {
	return e.Err.Error()
}

func (e Error) Marshal() ([]byte, error) {
	return json.MarshalIndent(e.JSON(), "", "  ")
}

func (e Error) JSON() interface{} {
	json := map[string]interface{}{
		"Error":  e.Err.Error(),
		"Type":   e.Type,
		"Values": e.Values,
	}
	return json
}
