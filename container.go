package medialocker

import (
	"github.com/codegangsta/inject"
)

type Container struct {
	inject.Injector
}

type ContainerConfigurator func(inject.Injector)

func NewContainer(configFns ...ContainerConfigurator) *Container {
	var i inject.Injector

	for x := len(configFns) - 1; x >= 0; x-- {
		injector := inject.New()
		fn := configFns[x]
		fn(injector)

		if i != nil {
			injector.SetParent(i)
		}

		i = injector
	}

	return &Container{Injector: i}
}
