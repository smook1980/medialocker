package medialocker

import (
	"testing"
	"github.com/codegangsta/inject"
)

func TestNewContainer(t *testing.T) {

	t.Run("works with no args", func(s *testing.T) {
		subject := NewContainer()
		if subject == nil {
			s.Error("Expected NewContainer to return *Container, got nil!")
		}
	})

	t.Run("it calls the config fn with the injector", func(s *testing.T) {
		wasCalledWithInjector := false

		configFn := func(injector inject.Injector) {
			if injector != nil {
				wasCalledWithInjector = true
			}
		}

		NewContainer(configFn)

		if !wasCalledWithInjector {
			s.Error("Expected config fn to be called with injector, instead was called nil.")
		}

	})
}

func TestContainer_Inject(t *testing.T) {
	
}

