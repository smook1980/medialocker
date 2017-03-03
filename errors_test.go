package medialocker

import (
	"errors"
	"fmt"
	"testing"
)

func TestMultiError_AddErrors(t *testing.T) {
	err1 := errors.New("Test Error 1")
	err2 := errors.New("Test Error 2")

	testTable := []struct {
		input, expected []error
		isNilReturned   bool
	}{
		{
			[]error{err1, err2},
			[]error{err1, err2},
			false,
		},
	}

	for x, data := range testTable {
		subject := MultiError(data.input...)
		if (subject == nil) != data.isNilReturned {
			msg := "Expected return type to "
			if !data.isNilReturned {
				msg = fmt.Sprintln(msg, "not")
			}

			msg = fmt.Sprintln(msg, " be nil!")
			t.Errorf("Test %d failed:\n%s", x, msg)
		}
	}

	// var err error = nil
	subject := MultiError(err1, err2)

	if subject == nil {
		t.Error("Expected error return type, got nil!")
	}

	mErr, ok := subject.(*Error)

	if !ok {
		t.Error("Expected returned error's underlying type to medialocker.Error!")
	}

	if len(mErr.errs) != 2 {
		t.Errorf("Expected two errors to have been record, got %d.", len(mErr.errs))
	}

	expectedErrs := []error{err1, err2}
	for x, err := range mErr.errs {
		if expectedErrs[x] != err {
			t.Errorf("Expected %d error to be %v.  Got: %v", x, expectedErrs[x], err)
		}
	}
}
