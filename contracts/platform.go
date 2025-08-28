package contracts

type Platform interface {
	Boot(container Container) (App, error)
}
