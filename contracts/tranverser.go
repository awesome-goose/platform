package contracts

type Transverser interface {
	Traverse(root Module) error
	Container() Container
	Routes() []Route
	OnBootHooks() Stack[func() error]
	OnShutdownHooks() Stack[func() error]
}
