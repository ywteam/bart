package ydk

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
var LogLevels = map[LogLevel][]string{
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

