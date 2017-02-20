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
	err := errs[0]
	errs = errs[1:]
	var mErr Error

	switch err := err.(type) {
	case *Error:
		err.addErrors(errs...)
		mErr = *err
	case Error:
		mErr = err
		mErr.addErrors(errs...)
	case error:
		mErr.addErrors(err)
		mErr.addErrors(errs...)
		//	case string:
		//		mErr.addErrors(errors.New(err))
		//		mErr.addErrors(errs...)
	case nil:
		mErr.addErrors(errs...)
	}

	return mErr
}

// Application Errors
var (
// ErrGeneric = errors.New("General locker error")
)

// func (e Error) Error() string {
// 	return string(e)
// }
