package argus

import (
	"encoding/json"
)

type ErrorType string

const (
	// database error types
	DbScanFailedError  ErrorType = "DbScanFailedError"
	DbCmdFailedError   ErrorType = "DbCmdFailedError"
	DbQueryFailedError ErrorType = "DbQueryFailedError"
	DbCloseFailedError ErrorType = "DbCloseFailedError"

	// io or json error types
	IOFailedReadError      ErrorType = "IOFailedReadError"
	IOFailedMarshalError   ErrorType = "IOFailedMarshalError"
	IOFailedUnmarshalError ErrorType = "IOFailedUnmarshalError"
)

type Error struct {
	Err    error
	Type   ErrorType
	Values map[string]interface{}
}

var Errors []Error

func (e Error) AppendTo(errs *[]Error) {
	*errs = append(*errs, e)
}

func (e *Error) SetValues(key string, value string) {
	e.Values[key] = value
}

func (e Error) String() string {
	return e.Err.Error()
}

func (e Error) Marshal() ([]byte, error) {
	return json.Marshal(e.JSON())
}

func (e Error) JSON() interface{} {
	json := map[string]interface{}{
		"Error":  e.Err.Error(),
		"Type":   e.Type,
		"Values": e.Values,
	}
	return json
}
