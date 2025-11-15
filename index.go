package platform

import "github.com/awesome-goose/contracts"

var (
	defaultKernel = core.NewKernel()
)

func Start(platform contracts.Platform, module contracts.Module) (func() error, error) {
	return defaultKernel.Start(platform, module)
}
