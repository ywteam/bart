package logger_test

import (
	"strings"
	"testing"

	logger "logger"
)
func TestLoggerMessage(t *testing.T) {
	message := logger.Message([]string{
		"red", "bold", "%s", "reset",
		"green", "italic", "%s", "reset",
		"blue", "dim", "%s", "reset",
		"yellow", "underline", "%s", "reset",
		"magenta", "reverse", "%s", "reset",
		"cyan", "bold", "%s", "reset",
	}, 
		"Hello, Red World!",
		"Hello, Green World!",
		"Hello, Blue World!",
		"Hello, Yellow World!",
		"Hello, Magenta World!",
		"Hello, Cyan World!",
	)
	if strings.Contains(message, "Hello") == false {
		t.Errorf("Expected Hello, got %s", message)
	}
	logger.Debug("Message: %s", message)
}
func TestLoggerMethods(t *testing.T) {
	logger.Debug("Hello, %s", "world")
	logger.Info("Hello, %s", "world")
	logger.Notice("Hello, %s", "world")
	logger.Warn("Hello, %s", "world")
	logger.Error("Hello, %s", "world")
	logger.Alert("Hello, %s", "world")
	logger.Critical("Hello, %s", "world")
	logger.Fatal("Hello, %s", "world")
}
func TestLogLevel(t *testing.T) {
	logLevels := logger.LogLevelsMap
	if len(logLevels) != 11 {
		t.Errorf("Expected 11, got %d", len(logLevels))
	}
	for level, logLevel := range logLevels {
		if logLevel.Priority < 0 {
			t.Errorf("Expected non-negative, got %d", logLevel.Priority)
		}
		if logLevel.Icon == "" {
			t.Errorf("Expected non-empty, got %s", logLevel.Icon)
		}
		if logLevel.Color == "" {
			t.Errorf("Expected non-empty, got %s", logLevel.Color)
		}
		logger.Debug("Level: %s, Priority: %d, Icon: %s, Color: %s", level, logLevel.Priority, logLevel.Icon, logLevel.Color)
	}
	
}
func TestNewContext(t *testing.T) {
	context := logger.NewContext("test", "development", "")
	if context.Env != "development" {
		t.Errorf("Expected development, got %s", context.Env)
	}
	if context.Pid == 0 {
		t.Errorf("Expected non-zero, got %d", context.Pid)
	}
	if context.Host == "" {
		t.Errorf("Expected non-empty, got %s", context.Host)
	}
	if context.Name != "test" {
		t.Errorf("Expected test, got %s", context.Name)
	}
	if context.StartAt.IsZero() {
		t.Errorf("Expected non-zero, got %s", context.StartAt)
	}
	if context.Format != "" {
		t.Errorf("Expected empty, got %s", context.Format)
	}
}