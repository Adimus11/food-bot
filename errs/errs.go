package errs

import "errors"

var (
	ErrNotFound      = errors.New("Resource not found")
	ErrTokenNotInCtx = errors.New("could not get token from ctx")
)
