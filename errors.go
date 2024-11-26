package flow

import "errors"

// Common errors
var (
	ErrProcessFailure            = errors.New("critical process failure")
	ErrParserFailure             = errors.New("critical parser failure")
	ErrSchemaVerificationFailure = errors.New("schema verification failed")
	ErrHandlerNotFound           = errors.New("process handler not found")
	ErrFlowNotFound              = errors.New("flow not defined in the schema")
	ErrInvalidRequest            = errors.New("invalid flow request")
	ErrInvalidRequestProcess     = errors.New("invalid flow request process")
	ErrInvalidRequestInputs      = errors.New("invalid flow request process defined inputs")
)
