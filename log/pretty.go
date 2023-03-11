package log

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type loggerImpl struct {
	Config
}

type Config struct {
	UseColors bool
	LogTime   bool
}

var prefixes = map[LogLevel]string{
	DebugLevel: "[DEBUG]",
	InfoLevel:  "[INFO] ",
	WarnLevel:  "[WARN] ",
	ErrorLevel: "[ERROR]",
}

var colorFuncs = map[LogLevel]func(string, ...interface{}) string{
	None:       color.WhiteString,
	DebugLevel: color.BlueString,
	InfoLevel:  color.GreenString,
	WarnLevel:  color.YellowString,
	ErrorLevel: color.RedString,
}

func (l loggerImpl) logString(level LogLevel, format string, args ...interface{}) string {
	var colorFunc func(string, ...interface{}) string
	if l.UseColors {
		colorFunc = colorFuncs[level]
	} else {
		colorFunc = colorFuncs[-1]
	}

	var timeString string
	if l.LogTime {
		timeString = fmt.Sprintf("%s - ", time.Now().Format("15:04:05"))
	}

	logString := fmt.Sprintf("%s %s%s", prefixes[level], timeString, fmt.Sprintf(format, args...))
	return colorFunc(logString)
}

func (l loggerImpl) log(level LogLevel, args ...interface{}) {
	l.logf(level, "%v", args...)
}

func (l loggerImpl) logf(level LogLevel, format string, args ...interface{}) {
	if dontLog(level) {
		return
	}
	fmt.Println(l.logString(level, format, args...))
}

func (l loggerImpl) Debug(args ...interface{}) {
	l.log(DebugLevel, args...)
}

func (l loggerImpl) Debugf(format string, args ...interface{}) {
	l.logf(DebugLevel, format, args...)
}

func (l loggerImpl) Info(args ...interface{}) {
	l.log(InfoLevel, args...)
}

func (l loggerImpl) Infof(format string, args ...interface{}) {
	l.logf(InfoLevel, format, args...)
}

func (l loggerImpl) Warn(args ...interface{}) {
	l.log(WarnLevel, args...)
}

func (l loggerImpl) Warnf(format string, args ...interface{}) {
	l.logf(WarnLevel, format, args...)
}

func (l loggerImpl) Error(args ...interface{}) {
	l.log(ErrorLevel, args...)
}

func (l loggerImpl) Errorf(format string, args ...interface{}) {
	l.logf(ErrorLevel, format, args...)
}
