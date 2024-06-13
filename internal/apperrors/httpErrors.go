package apperrors

import "net/http"

type HttpError struct {
	Message    string
	StatusCode int
}

func NewHttpError(message string, statusCode int) *HttpError {
	return &HttpError{
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e *HttpError) Error() string {
	return e.Message
}

var (
	ErrInternalServer = NewHttpError("Internal server error", http.StatusInternalServerError)
)
