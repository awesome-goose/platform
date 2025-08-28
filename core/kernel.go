package platform

import (
	"github.com/awesome-goose/platform/contracts"
)

type kernel struct {
	container   contracts.Container
	router      contracts.Router
	serializer  contracts.Serializer
	transverser contracts.Transverser
}

func NewKernel() *kernel {
	return &kernel{
		container:   NewContainer(),
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

	app, err := platform.Boot(k.container)
	if err != nil {
		return stop, err
	}

	err = app.Run(func(context contracts.Context) error {
		routes := k.transverser.Routes()
		segments := context.Segments()
		route, err := k.router.Find(routes, segments)
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

		// run route handler
		// use container to call the method on the struct, dynamically injecting params, queries, marshally body where nevessary
		output := route.Handler(context)
		err, ok := output.(error)
		if ok {
			return err
		}

		buf, err := k.serializer.Serialize(output)
		if err != nil {
			return err
		}

		err = context.Response().Write(buf)
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
