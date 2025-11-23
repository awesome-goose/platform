package env

import (
	"github.com/awesome-goose/contracts"
	"github.com/awesome-goose/platform/env/sources"
)

type Env struct {
	store map[string]string
}

func NewEnv() *Env {
	e := &Env{store: make(map[string]string)}
	e.FromSources(sources.NewOsEnvSource(), sources.NewFileEnvSource())
	return e
}

// FromSources applies multiple EnvSources to populate the env
func (e *Env) FromSources(sources ...contracts.EnvSource) {
	for _, src := range sources {
		src.Load(e)
	}
}

// Get returns a value from the store or the default if not found
func (e *Env) Get(key, defaultValue string) string {
	if val, ok := e.store[key]; ok {
		return val
	}
	return defaultValue
}

// Set assigns a value to the store
func (e *Env) Set(key, value string) {
	e.store[key] = value
}
