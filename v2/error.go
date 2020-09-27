package argus

// NewError creates new error.
func NewError(typ *ErrorType, err error) error {
	return Error{
		Type: typ,
		Err:  err,
	}
}

// Error is the base error type of this project
// which implements error interface.
type Error struct {
	// error type information
	Type *ErrorType

	// default error
	Err error

	// values for debugging
	Values map[string]interface{}
}

func (e Error) Error() string {
	return e.Err.Error()
}

// JSON parses Error to json format.
func (e Error) JSON() interface{} {
	if e.Err == nil {
		return nil
	}

	jsonMap := map[string]interface{}{
		"Type":   e.Type.Name,
		"Msg":    e.Type.Msg,
		"Error":  e.Err.Error(),
		"Values": e.Values,
	}
	return jsonMap
}

// ErrorType is the type of error.
type ErrorType struct {
	// error type
	Name string

	// errro messgae
	Msg string
}

// ErrorsHandler handles some errors occurred in the api stream.
type ErrorsHandler struct {
	Errs *[]error
}

// AppendError appends new error.
func (eh *ErrorsHandler) AppendError(err error) {
	*eh.Errs = append(*eh.Errs, err)
}

// Count returns the number of Errs elements.
func (eh ErrorsHandler) Count() int {
	return len(*eh.Errs)
}
