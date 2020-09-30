package argus

// Error is the base error type of this project
// which implements error interface.
type Error struct {
	// Type is error type information.
	Type *ErrorType

	// Err is default error.
	Err error

	// Values is for helping debug.
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
		"Pkg":    e.Type.Pkg,
		"Msg":    e.Type.Msg,
		"Error":  e.Err.Error(),
		"Values": e.Values,
	}
	return jsonMap
}

// NewError creates new error.
func NewError(typ *ErrorType, err error) error {
	return Error{
		Type: typ,
		Err:  err,
	}
}

// ErrorType is the type of error.
type ErrorType struct {
	// Pkg is the package name error is occurred.
	Pkg string

	// Msg is errro messgae.
	Msg string
}

// NewErrorType creates new error type.
func NewErrorType(pkg string, msg string) *ErrorType {
	return &ErrorType{
		Pkg: pkg,
		Msg: msg,
	}
}

// ErrorsHandler handles some errors occurred in the api call.
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
