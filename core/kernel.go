package core

import (
	"github.com/awesome-goose/contracts"
)

type kernel struct {
	router      contracts.RouterFinder
	serializer  contracts.Serializer
	transverser contracts.Transverser
}

func NewKernel() *kernel {
	return &kernel{
		router:      NewRouter(),
		serializer:  NewSerializer(),
		transverser: NewTransverser(),
	}
}

func (k *kernel) Start(platform contracts.Platform, module contracts.Module) (func() error, error) {
	stop := func() error {
		return k.transverser.OnShutdownHooks().ExecuteAll(func(fn func() error) error {
			return fn()
		})
	}

	err := k.transverser.Traverse(module)
	if err != nil {
		return stop, err
	}

	err = k.transverser.OnBootHooks().ExecuteAll(func(fn func() error) error {
		return fn()
	})
	if err != nil {
		return stop, err
	}

	app, err := platform.Boot(k.transverser.Container())
	if err != nil {
		return stop, err
	}

	err = app.Run(func(context contracts.Context) error {
		routes := k.transverser.Routes()
		route, err := k.router.Find(routes, context.Segments())
		if err != nil {
			return err
		}

		for _, middleware := range route.Middlewares {
			err := middleware.Handle(context)
			if err != nil {
				return err
			}
		}

		for _, validator := range route.Validators {
			err := validator.Validate(context)
			if err != nil {
				return err
			}
		}

		output := route.Handler(context)
		err, ok := output.(error)
		if ok {
			return err
		}

		serialType, buf, err := k.serializer.Serialize(output)
		if err != nil {
			return err
		}

		err = context.Response().Write(serialType, buf)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return stop, err
	}

	return stop, nil
}
