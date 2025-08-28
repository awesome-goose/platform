package contracts

type OnBoot interface {
	Boot() error
}

type OnShutdown interface {
	Shutdown() error
}
