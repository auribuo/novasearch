package log

func Setup(level string, noColor bool) error {
	SetDefault(New(&Config{
		UseColors: !noColor,
		LogTime:   true,
	}))
	return SetLogLevelFromString(level)
}
