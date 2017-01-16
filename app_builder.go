package medialocker

import (
	"os"
	"context"
	"sync"
)

type Injector func(*App) error

type AppBuilder struct {
	app *App
	errors []error
}

func NewAppBuilder() *AppBuilder {
	context, cancleFn := context.WithCancel(context.Background())

	return &AppBuilder{
		app: &App{
			In: os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
			context: context,
			wg: sync.WaitGroup{},
			cancleFn: cancleFn,
			Config: EmptyConfig,
			Registry: Registry{},
			Fs: LocalFileSystem(),
		},
	}
}

func (ab *AppBuilder) Build() (*App, []error) {
	if ab.app.Config == EmptyConfig {
		ab.app.Config, _ = BuildConfig()
	}

	ab.app.Log = NewLoggerWith(ab.app.Config)
	ab.app.Registry = *NewRegistry(ab.app.Log, ab.app.Config)

	return ab.app, ab.errors
}

func (ab *AppBuilder) Inject(injectors ...Injector) *AppBuilder {
	for _, fn := range(injectors) {
		err := fn(ab.app)
		if err != nil {
			ab.errors = append(ab.errors, err)
		}
	}

	return ab
}

func (ab *AppBuilder) WithConfiguration(configs ...Configuration) *AppBuilder {
	c, errs := BuildConfig(configs...)
	if len(errs) > 0 {
		ab.errors = append(ab.errors, errs...)
	}

	ab.app.Config = c

	return ab
}
