package errors

import "net/http"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	Status  int
	Code    string
	Message string
	Details []FieldError
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) StatusCode() int {
	return e.Status
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithCause(err error) *AppError {
	e.Err = err
	return e
}

func Internal(msg string) *AppError {
	if msg == "" {
		msg = "internal server error"
	}
	return &AppError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: msg}
}

func Unauthorized(msg string) *AppError {
	if msg == "" {
		msg = "You are not authenticated to access this resource."
	}
	return &AppError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: msg}
}

func NotFound(msg string) *AppError {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return &AppError{Status: http.StatusNotFound, Code: "NOT_FOUND", Message: msg}
}

func Forbidden(msg string) *AppError {
	if msg == "" {
		msg = "You are not authorized to perform the requested action."
	}
	return &AppError{Status: http.StatusForbidden, Code: "FORBIDDEN", Message: msg}
}

func BadRequest(msg string) *AppError {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return &AppError{Status: http.StatusBadRequest, Code: "BAD_REQUEST", Message: msg}
}

func Conflict(msg string) *AppError {
	if msg == "" {
		msg = "The resource already exists."
	}
	return &AppError{Status: http.StatusConflict, Code: "CONFLICT", Message: msg}
}

func Validation(msg string, details []FieldError) *AppError {
	if msg == "" {
		msg = "validation failed"
	}
	return &AppError{
		Status:  http.StatusBadRequest,
		Code:    "VALIDATION_ERROR",
		Message: msg,
		Details: details,
	}
}
