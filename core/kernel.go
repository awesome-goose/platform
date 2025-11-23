package core

import (
	"github.com/awesome-goose/contracts"
)

type kernel struct {
	router      contracts.RouterFinder
	routes      contracts.Routes
	serializer  contracts.Serializer
	transverser contracts.Transverser
}

func NewKernel() *kernel {
	return &kernel{
		router:      NewRouter(),
		routes:      []contracts.Route{},
		serializer:  NewSerializer(),
		transverser: NewTransverser(),
	}
}

func (k *kernel) Start(platform contracts.Platform, module contracts.Module) (func() error, error) {
	stop := func() error {
		return k.transverser.OnShutdownHooks().ExecuteAll(func(fn func(contracts.Kernel) error) error {
			return fn(k)
		})
	}

	container := k.transverser.Container()
	for _, fn := range services {
		container.Register(fn, "", true)
	}

	err := k.transverser.Traverse(module)
	if err != nil {
		return stop, err
	}

	err = k.transverser.OnBootHooks().ExecuteAll(func(fn func(contracts.Kernel) error) error {
		return fn(k)
	})
	if err != nil {
		return stop, err
	}

	app, err := platform.Boot(container)
	if err != nil {
		return stop, err
	}

	err = app.Run(func(context contracts.Context) error {
		routes := k.Routes()
		route, err := k.router.Find(routes, context.Request().Method().String(), context.Request().Paths())
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

func (k *kernel) Router() contracts.RouterFinder {
	return k.router
}

func (k *kernel) Routes() []contracts.Route {
	return k.routes
}

func (k *kernel) AppendRoutes(routes ...contracts.Route) []contracts.Route {
	k.routes = append(k.routes, routes...)
	return k.routes
}

func (k *kernel) Container() contracts.Container {
	return k.transverser.Container()
}
