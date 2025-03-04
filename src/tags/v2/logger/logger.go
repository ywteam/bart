package ydk

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
	console "github.com/ywteam/ydk-go/console"
)

const YDK_LOGGER_DEFAULT_CONTEXT = "🩳 YDK"

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
	State LoggerState
}
// SetContext sets the context of the YdkLogger.
// It takes a string parameter 'context' and updates the 'Context' field of the YdkLogger's state.
// It returns a pointer to the updated YdkLogger.
func (l *YdkLogger) SetContext(context string) *YdkLogger {
	// state.Context = context
	context = strings.ToUpper(context)
	message := console.Colorize(YDK_LOGGER_DEFAULT_CONTEXT, console.Yellow, console.Bold)
	if l.State.Context != context {
		// message = fmt.Sprintf("[%s/%s] ", YDK_LOGGER_DEFAULT_CONTEXT, context)
		message += "-" + console.Colorize(context, console.Gray, console.Italic)
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
	var leveRef = LogLevels[level]
	if len(leveRef) == 0 {
		leveRef = LogLevels[Info]
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
	// template = console.Colorize(template, console.Cyan, console.Normal)
	// log.Printf("%s | %s ", payload.Color, console.ColorIndex(payload.Color))
	data := map[string]interface{}{
		"Log": map[string]interface{}{
			"Team":      payload.Team,
			"Enabled":   payload.Enabled,
			"Env":       payload.Env,
			"Pid":       console.Colorize(fmt.Sprintf("%d", payload.Pid), console.Blue, console.Italic),
			"Host":      payload.Host,
			"IP":        payload.IP,
			"Context":   payload.Context,
			"Weigth":    payload.Weigth,
			"Level":     console.Colorize(payload.Level, console.ColorIndex(payload.Color), console.Underline),
			"Timestamp": payload.Timestamp,
			"Message":   payload.Message,
			"Elapsed":   console.Colorize(fmt.Sprintf("%.3f", payload.Elapsed), console.Gray, console.Italic),
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
	template = console.InterpolateEmojis(template)
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
}

var instance *YdkLogger

func Log(level LogLevel, message string) *YdkLogger {
	var text, payload = instance.privateWriteText(level, message)
	if payload.Enabled {
		text = instance.privateInterpolate(text, payload)
		log.Println(text)
	}
	return instance
}
// func Trace(message string) *YdkLogger {
// 	return Log(Trace, message)
// }
// func Debug(message string) *YdkLogger {
// 	return Log(Debug, message)
// }
// func Info(message string) *YdkLogger {
// 	return Log(Info, message)
// }
// func Notice(message string) *YdkLogger {
// 	return Log(Notice, message)
// }
// func Output(message string) *YdkLogger {
// 	return Log(Output, message)
// }
// func Success(message string) *YdkLogger {
// 	return Log(Success, message)
// }
// func Warning(message string) *YdkLogger {
// 	return Log(Warning, message)
// }
// func Alert(message string) *YdkLogger {
// 	return Log(Alert, message)
// }
// func Error(message string) *YdkLogger {
// 	return Log(Error, message)
// }
// func Critical(message string) *YdkLogger {
// 	return Log(Critical, message)
// }
// func Emergency(message string) *YdkLogger {
// 	return Log(Emergency, message)
// }
// func Panic(message string) *YdkLogger {
// 	return Log(Panic, message)
// }
// func Fatal(message string) *YdkLogger {
// 	return Log(Fatal, message)
// }

func init(){
	var once sync.Once
	once.Do(func() {
		instance = &YdkLogger{
			State: state,
		}
		instance.SetContext(state.Context)
	})
}
func Default() *YdkLogger {
	return instance
}
func NewLogger(context string) YdkLogger {
	state := Default().State
	return YdkLogger{
		State: LoggerState{
			Context: context,
			Env:     state.Env,
			Pid:     state.Pid,
			Host:    state.Host,
			IP:      state.IP,
		},
	}
}