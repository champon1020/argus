package argus

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

// JSON parses Error to json format.
func (e Error) JSON() interface{} {
	if e.Err == nil {
		return nil
	}

	jsonMap := map[string]interface{}{
		"Msg":    e.Msg,
		"Error":  e.Err.Error(),
		"Values": e.Values,
	}
	return jsonMap
}

// NewError creates new error.
func NewError(msg error, err error) *Error {
	return &Error{
		Msg: msg,
		Err: err,
	}
}
