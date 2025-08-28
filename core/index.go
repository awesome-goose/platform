package platform

import "github.com/awesome-goose/platform/contracts"

var (
	defaultKernel = NewKernel()
)

func Start(platform contracts.Platform, module contracts.Module) (func() error, error) {
	return defaultKernel.Start(platform, module)
}
