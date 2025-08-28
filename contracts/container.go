package contracts

type Container interface {
	Register(resolver any, name string, singleton bool) error
	Resolve(abstraction any, name string) error
	Make(value any) (any, error)
}
