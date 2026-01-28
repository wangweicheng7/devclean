package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Ignore []string `yaml:"ignore"`
	Rules  []string `yaml:"rules"`
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return &Config{}, nil
	}

	path := filepath.Join(home, ".devclean.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		// 文件不存在是正常情况
		return &Config{}, nil
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
