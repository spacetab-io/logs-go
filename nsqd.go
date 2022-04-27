package log

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

type NSQLogger struct {
	logger Logger
}

func NewNSQLogger(logger Logger) NSQLogger {
	return NSQLogger{logger: logger}
}

func (nl NSQLogger) Output(calldepth int, s string) error {
	output := strings.SplitN(s, " ", 2)

	if len(output) <= 1 {
		nl.defaultOutput(s)

		return nil
	}

	logLvl := parseNSQLogLvl(output[0])
	if logLvl == -2 {
		nl.defaultOutput(s)

		return nil
	}

	nl.logger.ForLogLevel(logLvl).Msg(strings.TrimSpace(output[1]))

	return nil
}

func (nl NSQLogger) defaultOutput(s string) {
	nl.logger.Info().Msg(s)
}

func parseNSQLogLvl(s string) zapcore.Level {
	var lvl zapcore.Level

	switch s {
	case "TRC", "DBG":
		lvl = zapcore.DebugLevel
	case "INF":
		lvl = zapcore.InfoLevel
	case "WRN":
		lvl = zapcore.WarnLevel
	case "ERR":
		lvl = zapcore.ErrorLevel
	case "FTL":
		lvl = zapcore.FatalLevel
	case "PNC":
		lvl = zapcore.PanicLevel
	default:
		lvl = -2
	}

	return lvl
}

func (nl NSQLogger) LogLevel() int {
	return nsqLogLvlFromZerologLogLvl(nl.logger.Level)
}

func nsqLogLvlFromZerologLogLvl(level zapcore.Level) int {
	var nsqlLL int

	switch level {
	case zapcore.DebugLevel:
		nsqlLL = 0
	case zapcore.InfoLevel:
		nsqlLL = 1
	case zapcore.WarnLevel:
		nsqlLL = 2
	case zapcore.ErrorLevel:
	case zapcore.FatalLevel:
	case zapcore.PanicLevel:
		nsqlLL = 3
	case -2:
		nsqlLL = 4
	}

	return nsqlLL
}
