package contracts

type Migration interface {
	Up() error
}
