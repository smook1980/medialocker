package medialocker

import (
	"context"
	"io"
	"sync"
)

// Set at build time, unless in DevMode
var (
	Version = ".dev"
	Commit  = ""
)

type Appliction interface {
	Shutdown()
}

type App struct {
	Out io.Writer
	In  io.Reader
	Err io.Writer
	Log *Logger

	Registry
	Config
	Fs FileSystem

	context  context.Context
	cancleFn context.CancelFunc
	wg       sync.WaitGroup
}

func (a *App) Shutdown() {
	wg := &a.wg
	registry := &a.Registry

	go func() {
		wg.Add(1)
		defer wg.Done()
		registry.Shutdown()
	}()

	a.cancleFn()
}

func (a *App) Wait() int {
	select {
	case <-a.context.Done():
	}

	(&a.wg).Wait()

	return a.ExitCode()
}

func (a *App) moduleStopped(name string, err error) {
	a.wg.Done()

	if err != nil {
		a.Log.Errorf("Module %s stopped unexpectedly with an error: %s", name, err)
	}
}

// Shutdown signals a task should finish and clean up
type ShutdownSignal interface{}

// Command is a blocking operation in the context of App
//
// *App is the application singleton
// <- ShutdownSignal signals command should abort.
//
// Returns error
type Command func(*App, chan<- ShutdownSignal) error

// type Module func(*App, chan<- ShutdownSignal) error
type Module func(*App) error

func (a *App) Start(name string, fn Module) {
	a.wg.Add(1)

	a.Log.Debugf("Starting %s module...", name)
	go func() {
		var err error
		defer a.moduleStopped(name, err)

		err = fn(a)
	}()
}

func (a *App) ExitCode() int {
	return 0
}
