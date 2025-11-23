package core

import (
	"github.com/awesome-goose/contracts"
	"github.com/awesome-goose/platform/errors"
)

type transverser struct {
	container       contracts.Container
	onBootStack     contracts.Stack[func(contracts.Kernel) error]
	onShutdownStack contracts.Stack[func(contracts.Kernel) error]
}

func NewTransverser() *transverser {
	return &transverser{
		container:       NewContainer(),
		onBootStack:     NewStack[func(contracts.Kernel) error](),
		onShutdownStack: NewStack[func(contracts.Kernel) error](),
	}
}

func (t *transverser) Traverse(root contracts.Module) error {
	return t.walk(root)
}

func (t *transverser) Container() contracts.Container {
	return t.container
}

func (t *transverser) OnBootHooks() contracts.Stack[func(contracts.Kernel) error] {
	return t.onBootStack
}

func (t *transverser) OnShutdownHooks() contracts.Stack[func(contracts.Kernel) error] {
	return t.onShutdownStack
}

func (t *transverser) walk(module contracts.Module) error {
	i, err := t.container.Make(module)
	if err != nil {
		return err
	}

	m, ok := i.(contracts.Module)
	if !ok {
		return errors.ErrInvalidModuleInstance
	}

	if bootable, ok := m.(contracts.Bootable); ok {
		t.onBootStack.Push(bootable.Boot)
	}

	if shutdownable, ok := m.(contracts.Shutdownable); ok {
		t.onShutdownStack.Push(shutdownable.Shutdown)
	}

	declarations := m.Declarations()

	for _, imp := range m.Imports() {
		declarations = append(declarations, imp.Exports()...)

		if err := t.walk(imp); err != nil {
			return err
		}
	}

	for _, declaration := range declarations {
		value, err := t.container.Make(declaration)
		if err != nil {
			return errors.ErrFailedToMakeDeclaration.WithError(err)
		}

		if bootable, ok := value.(contracts.Bootable); ok {
			t.onBootStack.Push(bootable.Boot)
		}

		if shutdownable, ok := value.(contracts.Shutdownable); ok {
			t.onShutdownStack.Push(shutdownable.Shutdown)
		}
	}

	return nil
}
