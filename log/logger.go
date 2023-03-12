package log

import "fmt"

var Default Logger = loggerImpl{
	Config: Config{
		UseColors: true,
		LogTime:   true,
	},
}

type LogLevel int

const (
	None LogLevel = iota - 1
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

var minLevel = InfoLevel

func SetLogLevel(level LogLevel) {
	minLevel = level
}

func dontLog(level LogLevel) bool {
	return level < minLevel
}

func SetLogLevelFromString(level string) error {
	switch level {
	case "debug":
		minLevel = DebugLevel
	case "info":
		minLevel = InfoLevel
	case "warn":
		minLevel = WarnLevel
	case "error":
		minLevel = ErrorLevel
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}
	return nil
}

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(err error)
	Errorf(format string, args ...interface{})
	log(level LogLevel, arg ...interface{})
	logf(level LogLevel, format string, arg ...interface{})
}

func SetDefault(l Logger) {
	Default = l
}

func New(cfg *Config) Logger {
	if cfg == nil {
		return Default
	}
	return loggerImpl{
		Config: *cfg,
	}
}
