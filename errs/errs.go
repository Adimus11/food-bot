package errs

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound                = errors.New("Resource not found")
	ErrTokenNotInCtx           = errors.New("could not get token from ctx")
	ErrWrongInterfaceAssertion = errors.New("wrong interfacer assertion")
)

func WrongInterfaceError(object interface{}, target string) error {
	return fmt.Errorf("%w: `%T` to `%s`", ErrWrongInterfaceAssertion, object, target)
}
