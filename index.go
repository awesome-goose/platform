package platform

import (
	"github.com/awesome-goose/contracts"
	"github.com/awesome-goose/platform/core"
)

var (
	defaultKernel = core.NewKernel()
)

func Start(platform contracts.Platform, module contracts.Module) (func() error, error) {
	return defaultKernel.Start(platform, module)
}
