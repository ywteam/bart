package ydk

// Scheme returns the current scheme for the logger.
const DEFAULT_TEMPLATE = "[{{.Log.Env}}] [{{.Log.Host}}@{{.Log.Pid}}] [{{.Log.Icon}} {{.Log.Level}}] {{.Log.Message}} [{{.Log.Elapsed}}] "

var SCHEMAS []LoggerScheme = []LoggerScheme{
	{Env: "development", Level: Trace, Template: DEFAULT_TEMPLATE},
	{Env: "production", Level: Alert, Template: DEFAULT_TEMPLATE},
}

// LoggerScheme represents the configuration scheme for the YdkLogger.
type LoggerScheme struct {
	Env      string   // Env represents the environment in which the logger is running.
	Level    LogLevel // Level represents the log level for the logger.
	Template string   // Template represents the log message template.
}

func (s *LoggerScheme) IsLevelEnabled(level LogLevel) bool {
	return level >= s.Level
}
func NewScheme(env string, level LogLevel, template string) {
	schema := LoggerScheme{Env: env, Level: level, Template: template}
	SCHEMAS = append(SCHEMAS, schema)
}
