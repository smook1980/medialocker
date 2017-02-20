package medialocker

import (
	"testing"

	"github.com/codegangsta/inject"
)

func TestNewContainer(t *testing.T) {
	ctx := NewTestAppCtx()
	t.Run("works with no args", func(s *testing.T) {
		subject := ctx.NewContainer()
		if subject == nil {
			s.Error("Expected NewContainer to return *Container, got nil!")
		}
	})

	t.Run("it calls the config fn with the injector", func(s *testing.T) {
		wasCalledWithInjector := false

		configFn := func(_ AppContext, injector inject.Injector) inject.Injector {
			if injector != nil {
				wasCalledWithInjector = true
			}

			return injector
		}

		ctx.NewContainer(configFn)

		if !wasCalledWithInjector {
			s.Error("Expected config fn to be called with injector, instead was called nil.")
		}

	})
}

func TestContainer_Inject(t *testing.T) {
	subject := NewTestAppCtx().NewContainer(AppContextConfig)
	obj := &struct {
		Logger *Logger `inject`
	}{}

	if err := subject.Apply(obj); err != nil {
		t.Errorf("Error returned injecting object: %s", err)
	}

	if obj.Logger == nil {
		t.Error("Container injected nil, expected Logger.")
	}

	obj.Logger.Error("Holy shit, it works.")
}
