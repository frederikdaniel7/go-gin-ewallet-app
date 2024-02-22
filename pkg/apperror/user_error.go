package apperror

type ErrorType struct {
	StatusCode int
	Message    string
}

func NewUserErrorType(statusCode int, message string) *ErrorType {
	return &ErrorType{
		StatusCode: statusCode,
		Message:    message,
	}
}
func NewInputErrorType(statusCode int, message string) *ErrorType {
	return &ErrorType{
		StatusCode: statusCode,
		Message:    message,
	}
}

func NewInternalErrorType(statusCode int, message string) *ErrorType {
	return &ErrorType{
		StatusCode: statusCode,
		Message:    message,
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
