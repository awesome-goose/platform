package config

import (
	"github.com/awesome-goose/utils/path"
)

var (
	Config = DefaultConfig()
)

func DefaultConfig() *config {
	configPath, err := path.Config()
	if err != nil {
		panic("failed to resolve config path: " + err.Error())
	}

	cfg, err := newConfig(configPath)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	return cfg
}
