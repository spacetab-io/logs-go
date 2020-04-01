package logs

import (
	"os"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

//NewLogger is logrus instantiating wrapper. Returns configured logrus instance
func NewLogger(config *Config) (log *Logger, err error) {
	log = &Logger{Logger: logrus.New()}
	log.Formatter = &logrus.TextFormatter{
		TimestampFormat:        time.RFC3339,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		QuoteEmptyFields:       true,
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.Out = os.Stdout

	// Flag for whether to log caller info (off by default)
	log.ReportCaller = true

	if log.Level, err = logrus.ParseLevel(config.LogLevel); err != nil {
		return nil, err
	}

	if config.Sentry != nil && !config.Sentry.Enable {
		return log, nil
	}

	if config.Stage == "" {
		config.SetStage()
	}

	if err := log.addSentryHook(config.Stage, config.Sentry); err != nil {
		return nil, err
	}

	return log, nil
}

func (l *Logger) addSentryHook(stage string, cfg *SentryConfig) error {
	sentryLevels := []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}

	hook, err := logrus_sentry.NewAsyncSentryHook(cfg.DSN, sentryLevels)
	if err != nil {
		return err
	}

	hook.SetEnvironment(stage)
	hook.StacktraceConfiguration.Enable = true
	hook.StacktraceConfiguration.Level = logrus.WarnLevel
	hook.StacktraceConfiguration.Skip = 6
	hook.StacktraceConfiguration.Context = 10
	hook.StacktraceConfiguration.IncludeErrorBreadcrumb = true
	hook.StacktraceConfiguration.SendExceptionType = true

	l.Hooks.Add(hook)

	return nil
}
