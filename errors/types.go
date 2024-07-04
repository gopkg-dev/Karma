package errors

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// statusReason returns the correct reason for the provided HTTP statuscode
func statusReason(status int) string {
	message := utils.StatusMessage(status)
	if message == "" {
		return "UnknownReason"
	}
	message = strings.TrimSpace(message)
	message = strings.ReplaceAll(message, "-", "")
	message = strings.ReplaceAll(message, " ", "")
	return message
}

// BadRequest generates a 400 error.
func BadRequest(format string, a ...interface{}) *Error {
	code := fiber.StatusBadRequest
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Reason:  statusReason(code),
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(format string, a ...interface{}) *Error {
	code := fiber.StatusUnauthorized
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Reason:  statusReason(code),
	}
}

// Forbidden generates a 403 error.
func Forbidden(format string, a ...interface{}) *Error {
	code := fiber.StatusForbidden
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Reason:  statusReason(code),
	}
}

// NotFound generates a 404 error.
func NotFound(format string, a ...interface{}) *Error {
	code := fiber.StatusNotFound
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Reason:  statusReason(code),
	}
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(format string, a ...interface{}) *Error {
	code := fiber.StatusMethodNotAllowed
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Reason:  statusReason(code),
	}
}

// TooManyRequests generates a 429 error.
func TooManyRequests(format string, a ...interface{}) *Error {
	code := fiber.StatusTooManyRequests
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Reason:  statusReason(code),
	}
}

// Timeout generates a 408 error.
func Timeout(format string, a ...interface{}) *Error {
	code := fiber.StatusRequestTimeout
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Reason:  statusReason(code),
	}
}

// Conflict generates a 409 error.
func Conflict(format string, a ...interface{}) *Error {
	code := fiber.StatusConflict
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Reason:  statusReason(code),
	}
}

// UnProcessableEntityError generates a 422 error.
func UnProcessableEntityError(format string, a ...interface{}) *Error {
	code := fiber.StatusUnprocessableEntity
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Reason:  statusReason(code),
	}
}

// UnsupportedMediaTypeError generates a 415 error.
func UnsupportedMediaTypeError(format string, a ...interface{}) *Error {
	code := fiber.StatusUnsupportedMediaType
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Reason:  statusReason(code),
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(format string, a ...interface{}) *Error {
	code := fiber.StatusInternalServerError
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Reason:  statusReason(code),
	}
}
