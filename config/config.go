package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type config struct {
	dir  string
	tree map[string]any
}

func newConfig(dir string) (*config, error) {
	tree := map[string]any{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		var data map[string]any
		if unmarshalErr := yaml.Unmarshal(content, &data); unmarshalErr != nil {
			return unmarshalErr
		}

		key := strings.TrimSuffix(filepath.Base(path), ext)
		tree[key] = data

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &config{
		tree: tree,
		dir:  dir,
	}, nil
}

func (c *config) Dir() string {
	return c.dir
}

func (c *config) Tree() map[string]any {
	return c.tree
}
