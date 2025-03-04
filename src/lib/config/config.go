package config

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
	"yellowteam/lib/logger"
)
type ConfigEntitry struct {
	Name     string
	Required bool
	Default  string
	Path     string	
}
func (c *ConfigEntitry) Get() (string, error) {
	if c.Path != "" {
		if secret, err := os.ReadFile(path.Join(c.Path, c.Name)); err == nil {
			return string(secret), nil
		}
	}
	if value := os.Getenv(c.Name); value != "" {
		return value, nil
	}
	if c.Default == "" && c.Required {
		return "", errors.New("Required environment variable not set: " + c.Name)
	}
	return c.Default, nil
}
func (c *ConfigEntitry) GetOrThrow() string {
	if value, err := c.Get(); err == nil {
		return value
	}
	panic("Required environment variable not set: " + c.Name)
}
func NewEnvVar(name string, required bool, defaultValue string) *ConfigEntitry {
	return &ConfigEntitry{
		Name: name,
		Required: required,
		Default: defaultValue,
	}
}
func NewSecret(path string, required bool, defaultValue string) *ConfigEntitry {
	name := path[strings.LastIndex(path, "/") + 1:]
	return &ConfigEntitry{
		Name: name,
		Path: path,
		Required: required,
		Default: defaultValue,
	}
}

type ConfigSection struct {
	Registry []ConfigEntitry
}
func (c *ConfigSection) IsValid() bool {
	isValid := true
	for _, entry := range c.Registry {
		if _, err := entry.Get(); err != nil && entry.Required {
			logger.Warn("Required config not set: " + entry.Name)
			isValid = false
		}
	}
	return isValid
}
func (c *ConfigSection) WithEnv(name string, required bool, defaultValue string) *ConfigSection {
	c.Registry = append(c.Registry, ConfigEntitry{
		Name: name,
		Required: required,
		Default: defaultValue,		
	})
	return c
}
func (c *ConfigSection) WithSecret(name string, path string, required bool, defaultValue string) *ConfigSection {
	c.Registry = append(c.Registry, ConfigEntitry{
		Name: name,
		Required: required,
		Default: defaultValue,
		Path: path,
	})
	return c
}
func (c *ConfigSection) Get(name string) string {
	for _, entry := range c.Registry {
		if entry.Name == name {
			if value, err := entry.Get(); err == nil {
				return value
			}
			logger.Warn("Required environment variable not set: " + name)
		}
	}
	return ""
}


type Config struct {
	Sections map[string]*ConfigSection
}
func (c *Config) Section(name string, registry ...ConfigEntitry) *ConfigSection {
	if c.Sections == nil {
		c.Sections = make(map[string]*ConfigSection)
	}
	if section, ok := c.Sections[name]; ok {
		registry = append(section.Registry, registry...)
	}
	c.Sections[name] = &ConfigSection{
		Registry: registry,
	}
	return c.Sections[name]
}

func (c *Config) WithDotEnv(dotenvs ...string) *Config {
	for _, dotenv := range dotenvs {
		if err := godotenv.Load(dotenv); err != nil {
			logger.Warn("Error loading .env file: " + dotenv)
		}
	}
	return c
}
func (c *Config) GetOrThrow(path string) string {
	if value := c.Get(path); value != "" {
		return value
	}
	panic("Required configuration not set: " + path)
}
func (c *Config) Get(path string) string {
	chunks := strings.Split(path, "/")
	if len(chunks) == 1 {
		if section, ok := c.Sections["/"]; ok {
			return section.Get(chunks[0])
		}
	} else if len(chunks) == 2 {
		if section, ok := c.Sections[chunks[0]]; ok {
			return section.Get(chunks[1])
		}
	}
	return ""
}
func (c *Config) IsValid() bool {
	isValid := true
	for name, section := range c.Sections {
		isValid = isValid && section.IsValid()
		if !isValid {
			logger.Warn("Invalid section configuartion: " + name)
		}
	}
	return isValid
}
func (c *Config) WithEnv(name string, required bool, defaultValue string) *Config {
	c.Section("/").WithEnv(name, required, defaultValue).IsValid()
	return c
}
func (c *Config) WithSecret(name string, path string, required bool, defaultValue string) *Config {
	c.Section("/").WithSecret(name, path, required, defaultValue).IsValid()
	return c
}

var instance *Config
func Default() *Config {
	if instance == nil {
		instance = &Config{}
	}
	return instance
}
func WithEnv(name string, required bool, defaultValue string) *Config {
	return Default().WithEnv(name, required, defaultValue)
}
func WithSecret(name string, path string, required bool, defaultValue string) *Config {
	return Default().WithSecret(name, path, required, defaultValue)
}
func WithDotEnv(dotenvs ...string) *Config {
	return Default().WithDotEnv(dotenvs...)
}
func Get(name string) string {
	return Default().Get(name)
}
func GetOrThrow(name string) string {
	return Default().GetOrThrow(name)
}
func GetEnv(name string) string {
	return os.Getenv(name)
}
func GetEnvOrDefault(name string, defaultValue string) string {
	if value := GetEnv(name); value != "" {
		return value
	}
	return defaultValue
}

