package log

import (
	"github.com/rs/zerolog"
)

type nsqLogger struct {
	logger zerolog.Logger
	level  zerolog.Level
}

func NSQLogger(logLevel string) nsqLogger {
	l := Logger()
	lvl, _ := zerolog.ParseLevel(logLevel)

	return nsqLogger{logger: l, level: lvl}
}

func (nl nsqLogger) Output(calldepth int, s string) error {
	WithLevel(nl.level).Msg(s)

	return nil
}
