package ydk

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	ydkConsole "github.com/ywteam/ydk-go/console"
)

const YDK_LOGGER_DEFAULT_CONTEXT = "🩳 YDK"

func GetRemoteAddr() string {
	// get remote addr from whatsmyip
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	if resp, err := httpClient.Get("https://api64.ipify.org?format=json"); err == nil {
		defer resp.Body.Close()
		if data, err := io.ReadAll(resp.Body); err == nil {
			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err == nil {
				if ip, ok := result["ip"]; ok {
					return ip.(string)
				}
			}
		}
	}
	return ""
}

// LogLevel represents the log levels for the YDK logger.
type LogLevel int

const (
	Trace     = iota // Trace represents
	Debug            // Debug represents
	Info             // Info represents
	Notice           // Notice represents
	Output           // Output represents
	Success          // Success represents
	Warning          // Warning represents
	Alert            // Alert represents
	Error            // Error represents
	Critical         // Critical represents
	Emergency        // Emergency represents
	Panic            // Panic represents
	Fatal            // Fatal represents
)

// LogLevel is a map that associates YdkLoggerLogLevel constants with their corresponding string representations.
var logLevels = map[LogLevel][]string{
	Trace:     {"TRACE", "🔍", "cyan"},
	Debug:     {"DEBUG", "🐞", "cyan"},
	Info:      {"INFO", "💬", "blue"},
	Notice:    {"NOTICE", "📢", "blue"},
	Output:    {"OUTPUT", "📌", "green"},
	Success:   {"SUCCESS", "👍", "green"},
	Warning:   {"WARNING", "⚠️ ", "yellow"},
	Alert:     {"ALERT", "🔔", "yellow"},
	Error:     {"ERROR", "⛔", "red", "red"},
	Critical:  {"CRITICAL", "🚨", "red"},
	Emergency: {"EMERGENCY", "🧯", "red"},
	Panic:     {"PANIC", "🔥", "red"},
	Fatal:     {"FATAL", "💀", "red"},
}

// Scheme represents the configuration scheme for the YdkLogger.
type Scheme struct {
	Env      string   // Env represents the environment in which the logger is running.
	Level    LogLevel // Level represents the log level for the logger.
	Template string   // Template represents the log message template.
}

func (s *Scheme) IsLevelEnabled(level LogLevel) bool {
	return level >= s.Level
}

// Scheme returns the current scheme for the logger.
const DEFAULT_TEMPLATE = "[{{.Log.Env}}] [{{.Log.Host}}@{{.Log.Pid}}] [{{.Log.Icon}} {{.Log.Level}}] {{.Log.Message}} [{{.Log.Elapsed}}] "

var SCHEMAS []Scheme = []Scheme{
	{Env: "development", Level: Trace, Template: DEFAULT_TEMPLATE},
	{Env: "production", Level: Alert, Template: DEFAULT_TEMPLATE},
}

func NewScheme(env string, level LogLevel, template string) {
	schema := Scheme{Env: env, Level: level, Template: template}
	SCHEMAS = append(SCHEMAS, schema)
}

// State represents the state of the logger.
type State struct {
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
func (s *State) Scheme() (Scheme, error) {
	for _, schema := range SCHEMAS {
		if schema.Env == s.Env {
			return schema, nil
		}
	}
	return Scheme{}, fmt.Errorf("no scheme found for environment %s", s.Env)
}

// Current logger state
var state = State{
	Env:     "development",
	Context: "logger",
	Pid:     os.Getpid(),
	Host: func() string {
		hostname, _ := os.Hostname()
		return hostname
	}(),
	StartAt: time.Now(),
	IP:      GetRemoteAddr(),
}

type YdkLoggerPayload struct {
	Team      string   `json:"team"`      // The team name of the logger.
	Enabled   bool     `json:"enabled"`   // Enabled represents whether the log level is enabled.
	Env       string   `json:"env"`       // The environment in which the logger is running.
	Pid       int      `json:"pid"`       // The process ID of the logger.
	Host      string   `json:"host"`      // The host name of the logger.
	IP        string   `json:"ip"`        // The IP address of the logger.
	Context   string   `json:"context"`   // The context of the logger.
	Weigth    LogLevel `json:"weight"`    // The weight of the log level.
	Level     string   `json:"level"`     // Level represents the log level for the logger.
	Timestamp int64    `json:"timestamp"` // Timestamp represents the time at which the log message was created.
	Message   string   `json:"message"`   // Message represents the log message.
	Elapsed   float64  `json:"elapsed"`   // Elapsed represents the time elapsed since the logger was started.
	Icon      string   `json:"icon"`      // Icon represents the icon associated with the log level.
	Color     string   `json:"color"`     // Color represents the color associated with the log level.
}

// YdkLogger represents a logger for the YDK application.
type YdkLogger struct {
	State State
}

// SetContext sets the context of the YdkLogger.
// It takes a string parameter 'context' and updates the 'Context' field of the YdkLogger's state.
// It returns a pointer to the updated YdkLogger.
func (l *YdkLogger) SetContext(context string) *YdkLogger {
	// state.Context = context
	context = strings.ToUpper(context)
	console := ydkConsole.YdkConsole{}
	message := console.Colorize(YDK_LOGGER_DEFAULT_CONTEXT, ydkConsole.Yellow, ydkConsole.Bold)
	if l.State.Context != context {
		// message = fmt.Sprintf("[%s/%s] ", YDK_LOGGER_DEFAULT_CONTEXT, context)
		message += "-" + console.Colorize(fmt.Sprintf("%s", context), ydkConsole.Gray, ydkConsole.Italic)
	}
	l.State.Context = context
	log.SetPrefix(fmt.Sprintf("[%s] ", message))

	// log.Default().SetPrefix(fmt.Sprintf("[%s]", YDK_LOGGER_DEFAULT_CONTEXT))
	return l
}

// IsLevelEnabled checks if the specified log level is enabled.
// It returns true if the log level is enabled, otherwise false.
func (l *YdkLogger) IsLevelEnabled(level LogLevel) bool {
	if scheme, err := l.State.Scheme(); err == nil {
		return scheme.IsLevelEnabled(level)
	}
	return false
}

func (l *YdkLogger) privateWriteJson(level LogLevel, message string) YdkLoggerPayload {
	var leveRef = logLevels[level]
	if len(leveRef) == 0 {
		leveRef = logLevels[Info]
	}
	var payload YdkLoggerPayload = YdkLoggerPayload{
		Team:      YDK_LOGGER_DEFAULT_CONTEXT,
		Context:   l.State.Context,
		Enabled:   l.IsLevelEnabled(level),
		Weigth:    level,
		Level:     leveRef[0],
		Icon:      leveRef[1],
		Color:     leveRef[2],
		Env:       l.State.Env,
		Pid:       l.State.Pid,
		Host:      l.State.Host,
		IP:        l.State.IP,
		Timestamp: time.Now().Unix(),
		Elapsed:   time.Since(l.State.StartAt).Seconds(),
		Message:   message,
	}
	return payload
}

func (l *YdkLogger) privateWriteText(level LogLevel, message string) (string, YdkLoggerPayload) {
	var payload = l.privateWriteJson(level, message)
	if scheme, err := l.State.Scheme(); err == nil {
		var template = scheme.Template
		var text = l.privateInterpolate(template, payload)
		// return fmt.Sprintf(template, payload), payload
		return fmt.Sprintf("%+v", text), payload
		// log.Printf("%+v", payload)
	}
	return message, payload
}

func (l *YdkLogger) privateInterpolate(template string, payload YdkLoggerPayload) string {
	envMap := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}
	console := ydkConsole.YdkConsole{}
	// template = console.Colorize(template, ydkConsole.Cyan, ydkConsole.Normal)
	// log.Printf("%s | %s ", payload.Color, console.ColorIndex(payload.Color))
	data := map[string]interface{}{
		"Log": map[string]interface{}{
			"Team":      payload.Team,
			"Enabled":   payload.Enabled,
			"Env":       payload.Env,
			"Pid":       console.Colorize(fmt.Sprintf("%d", payload.Pid), ydkConsole.Blue, ydkConsole.Italic),
			"Host":      payload.Host,
			"IP":        payload.IP,
			"Context":   payload.Context,
			"Weigth":    payload.Weigth,
			"Level":     console.Colorize(payload.Level, console.ColorIndex(payload.Color), ydkConsole.Underline),
			"Timestamp": payload.Timestamp,
			"Message":   payload.Message,
			"Elapsed":   console.Colorize(fmt.Sprintf("%.3f", payload.Elapsed), ydkConsole.Gray, ydkConsole.Italic),
			"Icon":      payload.Icon,
			"Color":     payload.Color,
		},
		"Env": envMap,
	}
	re := regexp.MustCompile(`{{*\.(.*?)}}`)
	matches := re.FindAllStringSubmatch(template, -1)
	for _, match := range matches {
		keys := strings.Split(match[1], ".")
		var current interface{} = data

		for _, key := range keys {
			if curMap, ok := current.(map[string]interface{}); ok {
				current, ok = curMap[key]
				if !ok {
					fmt.Printf("Key %s not found in map\n", key)
					break
				}
			} else {
				fmt.Printf("Current is not a map[string]interface{}\n")
				break
			}
		}

		if current != nil {
			template = strings.Replace(template, match[0], fmt.Sprintf("%v", current), -1)
		} else {
			fmt.Printf("Failed to convert value to string for key: %s, value: %v\n", match[1], current)
		}
		// if value, ok := current.(string); ok {
		// 	template = strings.Replace(template, match[0], value, -1)
		// } else {
		// 	fmt.Printf("Failed to convert value to string for key: %s, value: %v\n", match[1], current)
		// }
	}

	return template
}

// Log logs the given message at the specified log level.
// It checks if the log level is enabled and then prints the log message using the configured log scheme.
func (l *YdkLogger) Log(level LogLevel, message string) *YdkLogger {
	var text, payload = l.privateWriteText(level, message)
	if payload.Enabled {
		text = l.privateInterpolate(text, payload)
		log.Println(text)
	}
	return l
	// var payload = l.privateWriteJson(level, message)
	// if payload.Enabled {
	// 	// json, err := json.MarshalIndent(payload, "", "  ")
	// 	// json, err := json.Marshal(payload)
	// 	// if err != nil {
	// 	// 	log.Println(err)
	// 	// } else {
	// 	// 	log.Printf("%s", json)
	// 	// }
	// 	// if scheme, err := l.State.Scheme(); err == nil {
	// 	// 	var template = scheme.Template
	// 	// 	var logMessage = fmt.Sprintf(template, payload)
	// 	// 	// log.Printf("%+v", payload)
	// 	// 	log.Println(logMessage)
	// 	// } else {
	// 	// 	log.Println(message)
	// 	// }
	// 	// log.Writer().Write([]byte(fmt.Sprintf("%+v\n", payload)))
	// 	// log.Println(payload)

	// 	// log.Printf("%+v", payload)
	// 	log.Writer().Write([]byte(fmt.Sprintf("%+v\n", payload)))
	// }
	// return l
}
func (l *YdkLogger) Trace(message string) *YdkLogger {
	l.Log(Trace, message)
	return l
}
func (l *YdkLogger) Debug(message string) *YdkLogger {
	l.Log(Debug, message)
	return l
}
func (l *YdkLogger) Info(message string) *YdkLogger {
	l.Log(Info, message)
	return l
}
func (l *YdkLogger) Notice(message string) *YdkLogger {
	l.Log(Notice, message)
	return l
}
func (l *YdkLogger) Output(message string) *YdkLogger {
	l.Log(Output, message)
	return l
}
func (l *YdkLogger) Success(message string) *YdkLogger {
	l.Log(Success, message)
	return l
}
func (l *YdkLogger) Warning(message string) *YdkLogger {
	l.Log(Warning, message)
	return l
}
func (l *YdkLogger) Alert(message string) *YdkLogger {
	l.Log(Alert, message)
	return l
}
func (l *YdkLogger) Error(message string) *YdkLogger {
	l.Log(Error, message)
	return l
}
func (l *YdkLogger) Critical(message string) *YdkLogger {
	l.Log(Critical, message)
	return l
}
func (l *YdkLogger) Emergency(message string) *YdkLogger {
	l.Log(Emergency, message)
	return l
}
func (l *YdkLogger) Panic(message string) *YdkLogger {
	l.Log(Panic, message)
	return l
}
func (l *YdkLogger) Fatal(message string) *YdkLogger {
	l.Log(Fatal, message)
	return l
}

var instance *YdkLogger
var once sync.Once

func Default() *YdkLogger {
	once.Do(func() {
		instance = &YdkLogger{
			State: state,
		}
		instance.SetContext(state.Context)
	})
	return instance
}

func NewLogger(context string) YdkLogger {
	state := Default().State
	return YdkLogger{
		State: State{
			Context: context,
			Env:     state.Env,
			Pid:     state.Pid,
			Host:    state.Host,
			IP:      state.IP,
		},
	}
}
