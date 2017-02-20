package medialocker

import (
	"github.com/codegangsta/inject"
)

type Container struct {
	inject.Injector
}

type ContainerConfigurator func(AppContext, inject.Injector) inject.Injector

// NewContainer returns the application container configured to
// provide requested dependencies from the first ContainerConfiguration
// in the argument list able to satisfy the request.
//
// Example:
// (inject )
// testContext := NewContainer(
// 	RepoFakeWith(fakeDataForTesting),
// 	TestSettingsConfiguration(),
// 	TestLogSinkConfiguration(),
// 	AppConfiguration(),
// )
//
// result := MethodUnderTest(testContainer, blah)
func (c AppContext) NewContainer(configFns ...ContainerConfigurator) *Container {
	var i inject.Injector

	for x := len(configFns) - 1; x >= 0; x-- {
		injector := inject.New()
		fn := configFns[x]
		injector = fn(c, injector)

		if i != nil {
			injector.SetParent(i)
		}

		i = injector
	}

	return &Container{Injector: i}
}

func (c *Container) Inject(obj interface{}) error {
	return c.Apply(obj)
}
