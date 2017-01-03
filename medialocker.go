package medialocker

import (
	"log"
	"os"

	"golang.org/x/net/context"
)

// Application Errors
const (
	ErrGeneric = Error("General locker error")
)

type Error string

func (e Error) Error() string {
	return string(e)
}

// Logger represents an abstracted structured logging implementation. It
// provides methods to trigger log messages at various alert levels and a
// WithField method to set keys for a structured log message.
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})

	WithField(string, interface{}) Logger
}

type MediaLocker struct {
	context context.Context
}

func (ml *MediaLocker) Wait() {
	select {
	case <-ml.context.Done():
		return
	}
}

func (ml *MediaLocker) WaitAndExit() {
	select {
	case <-ml.context.Done():
		os.Exit(ml.ExitCode())
	}

	return
}

func (ml *MediaLocker) ExitCode() int {
	return 0
}

func AppInit(devMode bool) (*MediaLocker, error) {
	if devMode {
		log.Print("Dev Mode Enabled!")
	}

	return &MediaLocker{context: context.Background()}, nil
}
