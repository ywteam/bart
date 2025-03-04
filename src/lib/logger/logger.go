package logger

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"time"

	console "yellowteam/lib/console"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const YDK_LOGGER_DEFAULT_CONTEXT = "🩳" // 
// var loggerFormat = "%{context} %{icon} \t %{message} %{now} %{level} %{etime} %{Pid}"
type Context struct {
	Env     string    // The environment in which the logger is running.
	Pid     int       // The process ID of the logger.
	Host    string    // The host name of the logger.
	Name string    // The context of the logger.
	StartAt time.Time // The time at which the logger was started.
	Format  string    // The format of the logger.
	// IP      string    // The IP address of the logger.
	// User    string    // The user name of the logger.
}
func (s *Context) String() string {
	return fmt.Sprintf("Env: %s, Pid: %d, Host: %s, Context: %s, StartAt: %s", s.Env, s.Pid, s.Host, s.Name, s.StartAt)
}
func NewContext(name string, env string, format string) *Context {
	if env == "" { env = "development" }
	return &Context{
		Env:     env,
		Pid:     os.Getpid(),
		Host:    func() string { host, _ := os.Hostname(); return host }(),
		Name:    name,
		StartAt: time.Now(),
		Format:  format,
	}
}
type LogLevel struct {
	Priority int
	Icon string
	Color string
}
var LogLevelsMap = map[string]LogLevel{
	"fatal": {Priority: 0, Icon: "💀", Color: "red:bg"},
	"critical": {Priority: 1, Icon: "🔥", Color: "red"},
	"alert": {Priority: 2, Icon: "🚨", Color: "red"},
	"error": {Priority: 3, Icon: "💩", Color: "red"},
	"warning": {Priority: 4, Icon: "🔔", Color: "yellow"},
	"notice": {Priority: 5, Icon: "📌", Color: "blue"},
	"info": {Priority: 6, Icon: "📝", Color: "green"},
	"debug": {Priority: 7, Icon: "🐞", Color: "gray"},
	"trace": {Priority: 8, Icon: "🔍", Color: "gray"},
	"success": {Priority: 6, Icon: "👍", Color: "green"},
	"output": {Priority: 6, Icon: "📤", Color: "green"},
}

type Logger struct {
	// Context string
	Context *Context
	Console *console.Console
}
func (l *Logger) SetContext(context string) *Logger{
	l.Context.Name = context
	return l
}
func (l *Logger) SetFormat(format string) *Logger {
	l.Context.Format = format
	return l
}
func (l *Logger) ETime() string {
	return time.Since(l.Context.StartAt).String()
}
func (l *Logger) Log(level string, message string) *Logger {
	// apply Console.Colorize to message
	logLevel, ok := LogLevelsMap[level]
	if !ok { logLevel = LogLevelsMap["info"] }
	re := regexp.MustCompile(`%{(\w+)}`)
	template := l.Context.Format
	matches := re.FindAllStringSubmatch(template, -1)
	var metadata = map[string]interface{}{}
	for _, match := range matches {
		contextValue := reflect.ValueOf(l.Context).Elem().FieldByName(match[1])
		if contextValue.IsValid() {
			metadata[match[1]] = contextValue
			continue
		}
		switch match[1] {
		case "level":
			metadata["level"] = l.Console.Colorize(
				logLevel.Color,
				cases.Title(language.Und, cases.NoLower).String(level),
				"reset",
			)
		case "etime":
			metadata["etime"] = l.ETime()
		case "now":
			metadata["now"] = time.Now().Format("2006-01-02 15:04:05")
		case "icon":
			metadata["icon"] = logLevel.Icon
		case "context":
			metadata["context"] = l.Context.Name
		case "message":
			metadata["message"] = message
			// l.Console.Colorize(message...)
		}		
	}
	template = re.ReplaceAllStringFunc(template, func(s string) string {
		return fmt.Sprintf("%v", metadata[s[2:len(s)-1]])
	})
	l.Console.Println(template)
	return l
}
/**
 * Format message with fmt.Sprintf and apply Console.Colorize
 * @param format []string
 * @param a ...any
 * @return string
 * @example
 * 	logger.Message([]string{"red", "bold", "%s", "reset"}, "Hello, World!")
 */
func (l *Logger) Message(format []string, a ...any) string {
	return fmt.Sprintf(l.Console.Colorize(format...), a...)
}
func (l *Logger) Fatal(format string, a ...any) *Logger {
	return l.Log("fatal", fmt.Sprintf(format, a...))
}
func (l *Logger) Critical(format string, a ...any) *Logger {
	return l.Log("critical", fmt.Sprintf(format, a...))
}
func (l *Logger) Alert(format string, a ...any) *Logger {
	return l.Log("alert", fmt.Sprintf(format, a...))
}
func (l *Logger) Error(format string, a ...any) *Logger {
	return l.Log("error", fmt.Sprintf(format, a...))
}
func (l *Logger) Warn(format string, a ...any) *Logger {
	return l.Log("warning", fmt.Sprintf(format, a...))
}
func (l *Logger) Notice(format string, a ...any) *Logger {
	return l.Log("notice", fmt.Sprintf(format, a...))
}
func (l *Logger) Info(format string, a ...any) *Logger {
	return l.Log("info", fmt.Sprintf(format, a...))
}
func (l *Logger) Debug(format string, a ...any) *Logger {
	return l.Log("debug", fmt.Sprintf(format, a...))
}
func (l *Logger) Trace(format string, a ...any) *Logger {
	return l.Log("trace", fmt.Sprintf(format, a...))
}
func (l *Logger) Success(format string, a ...any) *Logger {
	return l.Log("success", fmt.Sprintf(format, a...))
}
func (l *Logger) Output(format string, a ...any) *Logger {
	return l.Log("output", fmt.Sprintf(format, a...))
}
func (l *Logger) Raize(format string, a ...any) error {
	return fmt.Errorf(l.Message([]string{"red", "bold", "%s", "reset"}, fmt.Sprintf(format, a...)))
}
func New(context string, term *console.Console) *Logger {
	if term == nil { term = console.Default() }
	return &Logger{
		Context: NewContext(
			context,
			os.Getenv("YDK_ENV"),
			"%{icon} %{message} \t %{level} %{etime} %{Pid} %{Env} %{context}",
		),
		Console: term,
	}
}

var instance *Logger
func Default() *Logger {
	if instance == nil { instance = New(YDK_LOGGER_DEFAULT_CONTEXT, console.Default()) }
	return instance
}

var Console = Default().Console
func Log(level string, message string) *Logger {
	return Default().Log(level, message)
}
func Message(format []string, a ...any) string {
	return Default().Message(format, a...)
}
func Fatal(format string, a ...any) *Logger {
	return Default().Fatal(format, a...)
}
func Critical(format string, a ...any) *Logger {
	return Default().Critical(format, a...)
}
func Alert(format string, a ...any) *Logger {
	return Default().Alert(format, a...)
}
func Error(format string, a ...any) *Logger {
	return Default().Error(format, a...)
}
func Warn(format string, a ...any) *Logger {
	return Default().Warn(format, a...)
}
func Notice(format string, a ...any) *Logger {
	return Default().Notice(format, a...)
}
func Info(format string, a ...any) *Logger {
	return Default().Info(format, a...)
}
func Debug(format string, a ...any) *Logger {
	return Default().Debug(format, a...)
}
func Trace(format string, a ...any) *Logger {
	return Default().Trace(format, a...)
}
func Success(format string, a ...any) *Logger {
	return Default().Success(format, a...)
}
func Output(format string, a ...any) *Logger {
	return Default().Output(format, a...)
}
func Raize(format string, a ...any) error {
	return Default().Raize(format, a...)
}
func SetContext(context string) *Logger {
	return Default().SetContext(context)
}
func Format(format string) *Logger {
	return Default().SetFormat(format)
}

// func init(){
// 	// _ = os.Setenv("NO_COLOR", "1")
// 	Log("info", "Message").Fatal("Message").Critical("Message").Alert("Message").Warn("Message").Notice("Message").Info("Message").Debug("Message").Trace("Message").Success("Message").Output("Message")
// 	Console.Write([]byte("Hello, World!"))
// }
