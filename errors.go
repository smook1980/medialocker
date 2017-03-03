package medialocker

import (
	"fmt"

	"github.com/pkg/errors"
)

// See https://github.com/pkg/errors
//
// Wrap:
// if err != nil {
// return errors.Wrap(err, "read failed")
// }
//
// Cause:
// switch err := errors.Cause(err).(type) {
// case *MyError:
// // handle specifically
// default:
// // unknown error
// }

type Error struct {
	errs []error
}

func WithStack(err error) error {
	return errors.WithStack(err)
}

func ErrCausedBy(err error, cause interface{}) (rerr error) {
	switch cause := cause.(type) {
	case string:
		rerr = errors.Wrap(err, cause)
	case error:
		rerr = errors.Wrap(err, cause.Error())
	}

	return rerr
}

func (me *Error) addErrors(err ...error) {
	me.errs = append(me.errs, err...)
}

func (me Error) Error() string {
	var errorStack string

	for n, err := range me.errs {
		errorStack = fmt.Sprintf("%s \n\nError %i:\n%s", errorStack, n, err)
	}

	return fmt.Sprintf("%i errors present:%s", len(me.errs), errorStack)
}

func MultiError(errs ...error) error {
	var merr *Error

	for _, err := range errs {
		if err == nil {
			continue
		}

		if merr == nil {
			switch err := err.(type) {
			case *Error:
				merr = err
			case Error:
				merr = &err
			case error:
				merr = &Error{}
				merr.addErrors(err)
			}
		} else {
			merr.addErrors(err)
		}
	}

	return merr
}
