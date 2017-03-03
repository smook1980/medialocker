package medialocker

import (
	"context"
	"io"
	"sync"

	"github.com/codegangsta/inject"
)

// Set at build time, unless in DevMode
var (
	Version = "devbuild"
	Commit  = ""
)

type appContextValueKey string

// AppContextConfig setup the default service providers
func AppContextConfig(ctx AppContext, i inject.Injector) inject.Injector {
	i.Map(ctx.Logger())
	i.Map(ctx.FileSystem())
	return i
}

const (
	FS_CTX_KEY  appContextValueKey = "fs"
	LOG_CTX_KEY appContextValueKey = "log"
)

type Appliction interface {
	Shutdown()
}

type AppContext struct {
	context.Context
}

func (ac AppContext) FileSystem() *FileSystem {
	switch fs := ac.Value(FS_CTX_KEY).(type) {
	case *FileSystem:
		return fs
	}

	return nil
}

func (ac AppContext) Logger() *Logger {
	return ac.Value(LOG_CTX_KEY).(*Logger)
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
	a.Wait()
	a.Log.WithField("prefix", "app").WithField("event", "shutdown").WithField("exit_code", a.ExitCode()).Infoln("Shutting down.")
}

func (a *App) Wait() int {
	(&a.wg).Wait()

	return 0
}

func (a *App) moduleStopped(name string, err error) {
	a.wg.Done()

	if err != nil {
		a.Log.WithField("prefix", name).WithField("event", "failed").Errorf("Stopped do to unexpected error %s", err)
	} else {
		a.Log.WithField("prefix", name).WithField("event", "finished").Infof("Module finished successfully.")
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
	a.Log.WithField("Version", Version).WithField("Head", Commit).Infoln("Media Locker is booting!")

	a.wg.Add(1)
	log := a.Log.WithField("prefix", name)

	log.Infoln("Starting module.")
	go func(a *App, log Log) {
		var err error

		err = fn(a)

		a.moduleStopped(name, err)
	}(a, log)
}

func (a *App) ExitCode() int {
	return 0
}
