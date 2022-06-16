package log

import (
	"go.uber.org/zap"
)

type GooseLogger struct {
	Logger
}

func NewGooseLogger(logger Logger) GooseLogger {
	l := logger
	l.Logger = l.Logger.WithOptions(zap.AddCallerSkip(1))

	return GooseLogger{Logger: l}
}

func (l GooseLogger) Fatal(v ...interface{}) {
	l.Logger.Fatal().Msgf("fatal error: %v", v)
}

func (l GooseLogger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatal().Msgf(format, v...)
}

func (l GooseLogger) Print(v ...interface{}) {
	l.Logger.Debug().Msgf("%v", v...)
}

func (l GooseLogger) Println(v ...interface{}) {
	l.Logger.Debug().Msgf("%v", v...)
}

func (l GooseLogger) Printf(format string, v ...interface{}) {
	l.Logger.Debug().Msgf(format, v...)
}
