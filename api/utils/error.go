package utils

import "fmt"

type ApiError struct {
	StatusCode int    `json:"-"`
	Reason     string `json:"error"`
}

func NewApiError(statusCode int, reason string) *ApiError {
	return &ApiError{
		StatusCode: statusCode,
		Reason:     reason,
	}
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("Could not finish request due: %s", e.Reason)
}
