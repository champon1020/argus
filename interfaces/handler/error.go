package handler

import "errors"

// Errors used as handler response.
var (
	ErrInvalidQueryParam  = errors.New("invalid query parameter")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrFailedDBExec       = errors.New("failed to execute the database command")
	ErrFailedGCSExec      = errors.New("failed to execute GCS api")
	ErrFailedOpenImage    = errors.New("cannot open the image as file")
)
