package response

type TError struct {
	StatusCode int
	CustomCode int
	Message    string
}

func (receiver *TError) Error() string {
	return receiver.Message
}

func NewError(statusCode, customCode int, message string) TError {
	return TError{
		StatusCode: statusCode,
		CustomCode: customCode,
		Message:    message,
	}
}
