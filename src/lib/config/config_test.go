package config_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"yellowteam/lib/config"
)

func TestConfigEntitry_Get(t *testing.T) {
	// Test with environment variable
	os.Setenv("TEST_ENV", "test_value")
	defer os.Unsetenv("TEST_ENV")

	entry := config.ConfigEntitry{Name: "TEST_ENV", Required: true}
	value, err := entry.Get()
	assert.NoError(t, err)
	assert.Equal(t, "test_value", value)

	// Test with default value
	entry = config.ConfigEntitry{Name: "NON_EXISTENT_ENV", Required: false, Default: "default_value"}
	value, err = entry.Get()
	assert.NoError(t, err)
	assert.Equal(t, "default_value", value)

	// Test with required but missing environment variable
	entry = config.ConfigEntitry{Name: "MISSING_ENV", Required: true}
	value, err = entry.Get()
	assert.Error(t, err)
	assert.Equal(t, "", value)

	// Test with secret file
	secretPath := path.Join(os.TempDir(), "secret.txt")
	os.WriteFile(secretPath, []byte("secret_value"), 0644)
	defer os.Remove(secretPath)

	entry = config.ConfigEntitry{Name: "secret.txt", Path: os.TempDir(), Required: true}
	value, err = entry.Get()
	assert.NoError(t, err)
	assert.Equal(t, "secret_value", value)
}

// func TestConfig_WithEnv(t *testing.T) {
// 	cfg := config.Default().WithEnv("TEST_ENV", true, "default_value")
// 	assert.Len(t, cfg.Registry, 1)
// 	assert.Equal(t, "TEST_ENV", cfg.Registry[0].Name)
// 	assert.Equal(t, true, cfg.Registry[0].Required)
// 	assert.Equal(t, "default_value", cfg.Registry[0].Default)
// }

// func TestConfig_WithSecret(t *testing.T) {
// 	cfg := config.Default().WithSecret("TEST_SECRET", "/path/to/secret", true, "default_value")
// 	assert.Len(t, cfg.Registry, 2)
// 	assert.Equal(t, "TEST_SECRET", cfg.Registry[1].Name)
// 	assert.Equal(t, "/path/to/secret", cfg.Registry[1].Path)
// 	assert.Equal(t, true, cfg.Registry[0].Required)
// 	assert.Equal(t, "default_value", cfg.Registry[0].Default)
// }

func TestConfig_WithDotEnv(t *testing.T) {
	cfg := config.Default().WithDotEnv(".env.test")
	assert.NotNil(t, cfg)
}

func TestConfig_GetOrThrow(t *testing.T) {
	config.Default().WithEnv("TEST_ENV", true, "")
	os.Setenv("TEST_ENV", "test_value")
	defer os.Unsetenv("TEST_ENV")

	value := config.GetOrThrow("TEST_ENV")
	assert.Equal(t, "test_value", value)

	assert.Panics(t, func() {
		config.GetOrThrow("MISSING_ENV")
	})
}

func TestConfig_Get(t *testing.T) {
	cfg := config.Default().WithEnv("TEST_ENV", true, "default_value")
	os.Setenv("TEST_ENV", "test_value")
	defer os.Unsetenv("TEST_ENV")

	value := cfg.Get("TEST_ENV")
	assert.Equal(t, "test_value", value)

	value = cfg.Get("MISSING_ENV")
	assert.Equal(t, "", value)
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_ENV", "test_value")
	defer os.Unsetenv("TEST_ENV")

	value := config.GetEnv("TEST_ENV")
	assert.Equal(t, "test_value", value)
}

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("TEST_ENV", "test_value")
	defer os.Unsetenv("TEST_ENV")

	value := config.GetEnvOrDefault("TEST_ENV", "default_value")
	assert.Equal(t, "test_value", value)

	value = config.GetEnvOrDefault("MISSING_ENV", "default_value")
	assert.Equal(t, "default_value", value)
}

func TestGet(t *testing.T) {
	config.Default().WithEnv("TEST_ENV", true, "default_value")
	os.Setenv("TEST_ENV", "test_value")
	defer os.Unsetenv("TEST_ENV")

	value := config.Get("TEST_ENV")
	assert.Equal(t, "test_value", value)
}

func TestGetOrThrow(t *testing.T) {
	cfg := config.Default().WithEnv("TEST_ENV", true, "")
	os.Setenv("TEST_ENV", "test_value")
	defer os.Unsetenv("TEST_ENV")

	value := cfg.GetOrThrow("TEST_ENV")
	assert.Equal(t, "test_value", value)

	assert.Panics(t, func() {
		config.GetOrThrow("MISSING_ENV")
	})
}
