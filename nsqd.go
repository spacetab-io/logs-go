package log

import (
	"strings"

	"github.com/rs/zerolog"
)

type NSQLogger struct {
	logger zerolog.Logger
}

func NewNSQLogger(logLevel string) NSQLogger {
	l := logger.With().CallerWithSkipFrameCount(4).Logger() //nolint:gomnd // Frame count to set correctly NSQ log line
	lvl, _ := zerolog.ParseLevel(logLevel)

	l.Level(lvl)

	return NSQLogger{logger: l}
}

func (nl NSQLogger) Output(calldepth int, s string) error {
	output := strings.SplitN(s, " ", 2)

	if len(output) <= 1 {
		nl.defaultOutput(s)

		return nil
	}

	logLvl := parseNSQLogLvl(output[0])
	if logLvl == zerolog.NoLevel {
		nl.defaultOutput(s)

		return nil
	}

	nl.logger.WithLevel(logLvl).Msg(strings.TrimSpace(output[1]))

	return nil
}

func (nl NSQLogger) defaultOutput(s string) {
	nl.logger.Log().Msg(s)
}

func parseNSQLogLvl(s string) zerolog.Level {
	var lvl zerolog.Level

	switch s {
	case "TRC":
		lvl = zerolog.TraceLevel
	case "DBG":
		lvl = zerolog.DebugLevel
	case "INF":
		lvl = zerolog.InfoLevel
	case "WRN":
		lvl = zerolog.WarnLevel
	case "ERR":
		lvl = zerolog.ErrorLevel
	case "FTL":
		lvl = zerolog.FatalLevel
	case "PNC":
		lvl = zerolog.PanicLevel
	default:
		lvl = zerolog.NoLevel
	}

	return lvl
}

func (nl NSQLogger) LogLevel() int {
	return nsqLogLvlFromZerologLogLvl(nl.logger.GetLevel())
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
