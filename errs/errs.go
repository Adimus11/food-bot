package errs

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound                = errors.New("Resource not found")
	ErrTokenNotInCtx           = errors.New("could not get token from ctx")
	ErrWrongInterfaceAssertion = errors.New("wrong interfacer assertion")
	ErrWrongMsgTypeInState     = errors.New("Couldn't this message type")
	ErrUnavailableTypeForUser  = errors.New("Unavailable event for user")
)

func WrongInterfaceError(object interface{}, target string) error {
	return fmt.Errorf("%w: `%T` to `%s`", ErrWrongInterfaceAssertion, object, target)
}
