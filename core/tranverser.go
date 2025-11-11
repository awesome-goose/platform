package platform

import (
	"fmt"

	"github.com/awesome-goose/contracts"
)

type transverser struct {
	container       contracts.Container
	routes          contracts.Routes
	onBootStack     contracts.Stack[func() error]
	onShutdownStack contracts.Stack[func() error]
}

func NewTransverser() *transverser {
	return &transverser{
		container:       NewContainer(),
		routes:          []contracts.Route{},
		onBootStack:     NewStack[func() error](),
		onShutdownStack: NewStack[func() error](),
	}
}

func (t *transverser) Traverse(root contracts.Module) error {
	return t.walk(root)
}

func (t *transverser) Container() contracts.Container {
	return t.container
}

func (t *transverser) Routes() []contracts.Route {
	return t.routes
}

func (t *transverser) OnBootHooks() contracts.Stack[func() error] {
	return t.onBootStack
}

func (t *transverser) OnShutdownHooks() contracts.Stack[func() error] {
	return t.onShutdownStack
}

func (t *transverser) walk(module contracts.Module) error {
	declarations := module.Declarations

	for _, imp := range module.Imports {
		declarations = append(declarations, imp.Exports...)

		if err := t.walk(imp); err != nil {
			return err
		}
	}

	for _, declaration := range declarations {
		value, err := t.container.Make(declaration)
		if err != nil {
			return fmt.Errorf("failed to make declaration %s: %w", declaration, err)
		}

		if route, ok := value.(contracts.Route); ok {
			t.routes = append(t.routes, route)
		}

		if bootable, ok := value.(contracts.OnBoot); ok {
			t.onBootStack.Push(bootable.Boot)
		}

		if shutdownable, ok := value.(contracts.OnShutdown); ok {
			t.onShutdownStack.Push(shutdownable.Shutdown)
		}
	}

	return nil
}
