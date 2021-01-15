package log

import (
	"github.com/rs/zerolog"
)

type NSQLogger struct {
	logger zerolog.Logger
	level  zerolog.Level
}

func NewNSQLogger(logLevel string) NSQLogger {
	l := Logger()
	lvl, _ := zerolog.ParseLevel(logLevel)

	return NSQLogger{logger: l, level: lvl}
}

func (nl NSQLogger) Output(calldepth int, s string) error {
	WithLevel(nl.level).Msg(s)

	return nil
}
