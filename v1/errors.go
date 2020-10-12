package argus

import (
	"encoding/json"
	"errors"
)

type ErrorType string

const (
	// basic error
	BasicError ErrorType = "BasicError"

	// database error types
	DbRuntimeError     ErrorType = "DbRuntimeError"
	DbScanFailedError  ErrorType = "DbScanFailedError"
	DbCmdFailedError   ErrorType = "DbCmdFailedError"
	DbQueryFailedError ErrorType = "DbQueryFailedError"
	DbCloseFailedError ErrorType = "DbCloseFailedError"

	// io or json error types
	IOFailedOpenError      ErrorType = "IOFailedOpenError"
	IOFailedCloseError     ErrorType = "IOFailedCloseError"
	IOFailedReadError      ErrorType = "IOFailedReadError"
	IOFailedWriteError     ErrorType = "IOFailedWriteError"
	IOFailedRemoveError    ErrorType = "IOFailedRemoveError"
	IOFailedMarshalError   ErrorType = "IOFailedMarshalError"
	IOFailedUnmarshalError ErrorType = "IOFailedUnmarshalError"

	// time error types
	TimeFailedParseError ErrorType = "TimeFailedParseError"

	// argus pkg errors
	ConfigFailedLoadError ErrorType = "ConfigFailedLoadError"

	// multi format data
	MultiFormatFailedOpenError ErrorType = "MultiFormatFailedOpenError"

	// authenticate
	AuthFailedVerifyError ErrorType = "AuthFailedVerifyError"
	JwtFailedParseError   ErrorType = "JwtFailedParseError"
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
	return json.Marshal(e.JSON())
}

func (e Error) MarshalIndent() ([]byte, error) {
	return json.MarshalIndent(e.JSON(), "", "  ")
}

func (e Error) JSON() interface{} {
	if e.Err == nil {
		e.Err = errors.New("empty")
	}

	jsonMap := map[string]interface{}{
		"Error":  e.Err.Error(),
		"Type":   e.Type,
		"Values": e.Values,
	}
	return jsonMap
}
