package contracts

type App interface {
	Run(fn func(c Context) error) error
}
