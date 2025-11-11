package env

import "github.com/awesome-goose/contracts"

type env struct {
	store map[string]string
}

func newEnv() *env {
	return &env{store: make(map[string]string)}
}

// FromSources applies multiple EnvSources to populate the env
func (e *env) FromSources(sources ...contracts.EnvSource) {
	for _, src := range sources {
		src.Load(e)
	}
}

// Get returns a value from the store or the default if not found
func (e *env) Get(key, defaultValue string) string {
	if val, ok := e.store[key]; ok {
		return val
	}
	return defaultValue
}

// Set assigns a value to the store
func (e *env) Set(key, value string) {
	e.store[key] = value
}
