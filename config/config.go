package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/awesome-goose/utils/path"
	"gopkg.in/yaml.v3"
)

type AppConfigPath string

func (p AppConfigPath) String() string {
	return string(p)
}

type Config struct {
	dir  string
	tree map[string]any
}

func NewConfig(appPath AppConfigPath) (*Config, error) {
	dir := appPath.String()
	if dir == "" {
		defaultDir, err := path.Config()
		if err != nil {
			panic("failed to resolve config path: " + err.Error())
		}

		dir = defaultDir
	}

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

	return &Config{
		tree: tree,
		dir:  dir,
	}, nil
}

func (c *Config) Dir() string {
	return c.dir
}

func (c *Config) Tree() map[string]any {
	return c.tree
}
