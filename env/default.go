package env

import "github.com/awesome-goose/platform/env/sources"

var (
	Env = defaultEnv()
)

func defaultEnv() *env {
	e := newEnv()
	e.FromSources(sources.NewOsEnvSource(), sources.NewFileEnvSource())
	return e
}
