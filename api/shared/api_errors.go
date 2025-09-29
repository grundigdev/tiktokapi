package shared

type NotFoundError struct {
	Message string
}

func (receiver NotFoundError) Error() string {
	return receiver.Message
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{
		Message: message,
	}
}
