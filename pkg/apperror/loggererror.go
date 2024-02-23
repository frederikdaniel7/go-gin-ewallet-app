package apperror

import "runtime/debug"

type AppError struct {
	// other various data for your errors
	stack []byte
}

func NewAppError() *AppError {
	return &AppError{
		stack: debug.Stack(),
	}
}

func (e *AppError) Error() string {
	return "got error!"
}


