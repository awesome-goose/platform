package config

import (
	"fmt"
	"strings"
)

// Get retrieves a value using dotted path (e.g., "db.host") and renders {{ENV_VAR}}.
func (c *config) Get(path string) (any, error) {
	parts := strings.Split(path, ".")
	current := any(c.tree)
	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("path %s not found (at %s)", path, part)
		}
		current, ok = m[part]
		if !ok {
			return nil, fmt.Errorf("key '%s' not found in path '%s'", part, path)
		}
	}
	if str, ok := current.(string); ok {
		return renderEnv(str), nil
	}
	return current, nil
}

// Set sets a value using dotted path (e.g., "db.host").
func (c *config) Set(path string, value any) error {
	parts := strings.Split(path, ".")
	lastKey := parts[len(parts)-1]
	current := c.tree
	for _, part := range parts[:len(parts)-1] {
		next, ok := current[part]
		if !ok {
			next = map[string]any{}
			current[part] = next
		}
		asMap, ok := next.(map[string]any)
		if !ok {
			return fmt.Errorf("cannot set key on non-map value at '%s'", part)
		}
		current = asMap
	}
	current[lastKey] = value
	return nil
}
