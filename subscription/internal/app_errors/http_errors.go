package apperrors

import "net/http"

type HttpError struct {
	Message      string
	StatusCode   int
	JSONResponse map[string]interface{}
}

func NewHttpError(message string, statusCode int, jsonResponse map[string]interface{}) *HttpError {
	return &HttpError{
		Message:      message,
		StatusCode:   statusCode,
		JSONResponse: jsonResponse,
	}
}

func (e *HttpError) Error() string {
	return e.Message
}

var (
	ErrInternalServer = NewHttpError("Internal server error",
		http.StatusInternalServerError,
		map[string]interface{}{"error": "internal server error"})
	ErrSubscriptionAlreadyExists = NewHttpError("Already Exists",
		http.StatusBadRequest,
		map[string]interface{}{
			"error": "subscription already exists",
		})
	ErrDatabase = NewHttpError("Database error",
		http.StatusInternalServerError,
		map[string]interface{}{
			"error": "database raised an error",
		})
	ErrAlreadyUnsubscribed = NewHttpError("subscription is already deactivated", http.StatusInternalServerError, map[string]interface{}{
		"error": "subscription is already deactivated",
	})
)
