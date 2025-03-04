package ydk

import (
	"fmt"
	"os"
	"time"
)

// LoggerState represents the state of the logger.
type LoggerState struct {
	Env     string    // The environment in which the logger is running.
	Pid     int       // The process ID of the logger.
	Host    string    // The host name of the logger.
	Context string    // The context of the logger.
	StartAt time.Time // The time at which the logger was started.
	IP      string    // The IP address of the logger.
}

// Scheme returns the Scheme associated with the current State's environment.
// It iterates through the list of Schemas and returns the first one that matches the current environment.
// If no matching Scheme is found, it returns an empty Scheme and an error.
func (s *LoggerState) Scheme() (LoggerScheme, error) {
	for _, schema := range SCHEMAS {
		if schema.Env == s.Env {
			return schema, nil
		}
	}
	return LoggerScheme{}, fmt.Errorf("no scheme found for environment %s", s.Env)
}

// Current logger state
var state = LoggerState{
	Env:     "development",
	Context: "logger",
	Pid:     os.Getpid(),
	Host: func() string {
		hostname, _ := os.Hostname()
		return hostname
	}(),
	StartAt: time.Now(),
	IP:      "GetRemoteAddr()",
}
