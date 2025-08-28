package contracts

type Kernel interface {
	Start(platform Platform, module Module) (stop func() error, err error)
}
