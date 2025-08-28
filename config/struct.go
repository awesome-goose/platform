package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ToStruct loads a specific YAML file from disk and unmarshals it into the given struct.
func (c *config) Export(namespace string, config any) error {
	// Try both .yaml and .yml extensions
	extensions := []string{".yaml", ".yml"}

	var filePath string
	found := false
	for _, ext := range extensions {
		path := filepath.Join(c.Dir(), namespace+ext)
		if _, err := os.Stat(path); err == nil {
			filePath = path
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("config file for namespace '%s' not found", namespace)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(content, config); err != nil {
		return fmt.Errorf("failed to unmarshal config to struct: %w", err)
	}

	return nil
}

func (c *config) Import(namespace string, in any) error {
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	// Marshal the struct to JSON first
	bytes, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	// Unmarshal back to map[string]any
	var asMap map[string]any
	if err := json.Unmarshal(bytes, &asMap); err != nil {
		return fmt.Errorf("failed to convert struct to map: %w", err)
	}

	// Store in the internal tree
	c.tree[namespace] = asMap
	return nil
}
