package log

import (
	"github.com/rs/zerolog"
)

type NSQLogger struct {
	level zerolog.Level
}

func NewNSQLogger(logLevel string) NSQLogger {
	lvl, _ := zerolog.ParseLevel(logLevel)

	return NSQLogger{level: lvl}
}

func (nl NSQLogger) Output(calldepth int, s string) error {
	WithLevel(nl.level).Msg(s)

	return nil
}

func (nl NSQLogger) LogLevel() int {
	return nsqLogLvlFromZerologLogLvl(nl.level)
}

func nsqLogLvlFromZerologLogLvl(level zerolog.Level) int {
	var nsqlLL int

	switch level {
	case zerolog.TraceLevel:
	case zerolog.DebugLevel:
		nsqlLL = 0
	case zerolog.InfoLevel:
		nsqlLL = 1
	case zerolog.WarnLevel:
		nsqlLL = 2
	case zerolog.ErrorLevel:
	case zerolog.FatalLevel:
	case zerolog.PanicLevel:
		nsqlLL = 3
	case zerolog.NoLevel:
	case zerolog.Disabled:
		nsqlLL = 4
	}

	return nsqlLL
}
