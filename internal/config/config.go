package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Project ProjectConfig `yaml:"project"`
	Env     EnvConfig     `yaml:"env"`
	Compose ComposeConfig `yaml:"compose"`
	Checks  ChecksConfig  `yaml:"checks"`
}

type ProjectConfig struct {
	Name string `yaml:"name"`
}

type EnvConfig struct {
	File     string   `yaml:"file"`
	Required []string `yaml:"required"`
}

type ComposeConfig struct {
	ProjectName string   `yaml:"project_name"`
	Files       []string `yaml:"files"`
}

type ChecksConfig struct {
	Commands []string `yaml:"commands"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg.ApplyDefaults()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Project.Name) == "" {
		return fmt.Errorf("project.name is required")
	}

	if len(c.Compose.Files) == 0 {
		return fmt.Errorf("compose.files is required")
	}

	return nil
}

func (c *Config) ApplyDefaults() {
	if strings.TrimSpace(c.Env.File) == "" {
		c.Env.File = ".env"
	}

	if strings.TrimSpace(c.Compose.ProjectName) == "" {
		c.Compose.ProjectName = c.Project.Name
	}

	if len(c.Checks.Commands) == 0 {
		c.Checks.Commands = []string{"docker"}
	}
}
