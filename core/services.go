package core

import (
	"github.com/awesome-goose/contracts"
	"github.com/awesome-goose/platform/config"
	"github.com/awesome-goose/platform/env"
	"github.com/awesome-goose/platform/log"
	"github.com/awesome-goose/platform/log/formatters"
	"github.com/awesome-goose/platform/log/modifiers"
	"github.com/awesome-goose/platform/log/processors"
	"github.com/awesome-goose/utils/path"
)

var (
	services = []any{
		func() (config.AppConfigPath, error) {
			path, err := path.AppRoot()
			return config.AppConfigPath(path), err
		},
		func() (log.AppLogChannel, error) {
			return log.AppLogChannel("std"), nil
		},
		func() (log.AppLoggers, error) {
			return []*log.Logger{
				log.NewLogger(
					[]contracts.Modifier{
						modifiers.NewUUID(),
						modifiers.NewColorTagsModifier(),
						modifiers.NewSystemInfo(),
						modifiers.NewStackTrace(),
					},
					formatters.NewJSON(),
					processors.NewConsole(),
				),
			}, nil
		},
		func(path config.AppConfigPath) (contracts.Config, error) {
			return config.NewConfig(path)
		},
		func(path config.AppConfigPath) (contracts.Env, error) {
			return env.NewEnv(), nil
		},
		func(channel log.AppLogChannel, loggers []*log.Logger) (contracts.Log, error) {
			return log.NewLog(channel, loggers...), nil
		},
	}
)
