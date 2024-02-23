package apperror

type ErrorType struct {
	StatusCode int
	Message    string
	Stack      []byte
}

func NewUserErrorType(statusCode int, message string, stack []byte) *ErrorType {
	return &ErrorType{
		StatusCode: statusCode,
		Message:    message,
	}
}
func NewInputErrorType(statusCode int, message string, stack []byte) *ErrorType {
	return &ErrorType{
		StatusCode: statusCode,
		Message:    message,
		Stack:      stack,
	}
}

func NewInternalErrorType(statusCode int, message string, stack []byte) *ErrorType {
	return &ErrorType{
		StatusCode: statusCode,
		Message:    message,
		Stack:      stack,
	}
}

func NewCredentialsErrorType(statusCode int, message string) *ErrorType {
	return &ErrorType{
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e *ErrorType) Error() string {
	return e.Message
}

func (e *ErrorType) GetStackTrace() []byte {
	return e.Stack
}
