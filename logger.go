package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/getsentry/sentry-go"
	"github.com/spacetab-io/configuration-structs-go/v2/contracts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ErrEmptyOutput = errors.New("log output is not defined")

func contextFields(ctx context.Context) (fields map[contracts.ContextKey]interface{}) {
	fields = make(map[contracts.ContextKey]interface{})
	if requestID, ok := ctx.Value(contracts.ContextKeyRequestID).(fmt.Stringer); ok &&
		requestID.String() != "00000000-0000-0000-0000-000000000000" {
		fields[contracts.ContextKeyRequestID] = requestID.String()
	}

	if requestID, ok := ctx.Value(contracts.ContextKeyRequestID).(string); ok &&
		requestID != "00000000-0000-0000-0000-000000000000" {
		fields[contracts.ContextKeyRequestID] = requestID
	}

	return fields
}

type Logger struct {
	*zap.Logger
	Level zapcore.Level
	cfg   contracts.LogsCfgInterface
}

func Init(cfg contracts.LogsCfgInterface, stage string, serviceAlias string, serviceVersion string, w io.Writer) (Logger, error) {
	logLevel, err := zapcore.ParseLevel(cfg.GetLevel())
	if err != nil {
		return Logger{}, fmt.Errorf("error level parsing error: %w", err)
	}

	if w == nil {
		return Logger{}, ErrEmptyOutput
	}

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < logLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.AddSync(w)
	consoleErrors := zapcore.AddSync(w)

	logConfig := zap.NewProductionEncoderConfig()
	logConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	logConfig.EncodeCaller = zapcore.FullCallerEncoder
	logConfig.ConsoleSeparator = " | "
	logConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if cfg.IsColored() {
		logConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var enc zapcore.Encoder

	switch cfg.GetFormat() {
	case "json":
		enc = zapcore.NewJSONEncoder(logConfig)
	case "text":
		enc = zapcore.NewConsoleEncoder(logConfig)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(enc, consoleErrors, highPriority),
		zapcore.NewCore(enc, consoleDebugging, lowPriority),
	}

	if cfg.IsSentryEnabled() {
		sentryClient, err := sentry.NewClient(sentry.ClientOptions{
			Dsn:   cfg.GetSentryDSN(),
			Debug: cfg.SentryDebugEnabled(),
			// AttachStacktrace: true,
			ServerName:  serviceAlias,
			Release:     serviceVersion,
			Environment: stage,
			Integrations: func(integrations []sentry.Integration) []sentry.Integration {
				filtered := integrations[:0]
				for _, integration := range integrations {
					if integration.Name() != "Modules" {
						filtered = append(filtered, integration)
					}
				}

				return filtered
			},
		})
		if err != nil {
			return Logger{}, fmt.Errorf("sentry client init error: %w", err)
		}

		buf := &bytes.Buffer{}
		sentryErrors := zapcore.AddSync(&zapcore.BufferedWriteSyncer{WS: zapcore.AddSync(buf)})
		cores = append(
			cores,
			zapcore.NewCore(
				&sentryEncoder{Encoder: zapcore.NewJSONEncoder(logConfig), client: sentryClient}, sentryErrors, highPriority,
			),
		)
	}

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	coresTee := zapcore.NewTee(cores...)

	logger := zap.New(coresTee, zap.WithCaller(cfg.ShowCaller()), zap.AddCallerSkip(cfg.GetCallerSkipFrames()))

	defer func() {
		_ = logger.Sync()
	}()

	zl := Logger{Logger: logger, Level: logLevel, cfg: cfg}

	return zl, nil
}

type Event struct {
	*Logger
	callerSkip int
	lvl        zapcore.Level
	fields     []zapcore.Field
}

func (l Logger) newEvent(lvl zapcore.Level) *Event {
	return &Event{Logger: &l, lvl: lvl, callerSkip: l.cfg.GetCallerSkipFrames()}
}
func (l Logger) Debug() *Event                        { return l.newEvent(zapcore.DebugLevel) }
func (l Logger) Info() *Event                         { return l.newEvent(zapcore.InfoLevel) }
func (l Logger) Warn() *Event                         { return l.newEvent(zapcore.WarnLevel) }
func (l Logger) Error() *Event                        { return l.newEvent(zapcore.ErrorLevel) }
func (l Logger) DPanic() *Event                       { return l.newEvent(zapcore.DPanicLevel) }
func (l Logger) Panic() *Event                        { return l.newEvent(zapcore.PanicLevel) }
func (l Logger) Fatal() *Event                        { return l.newEvent(zapcore.FatalLevel) }
func (l Logger) ForLogLevel(lvl zapcore.Level) *Event { return l.newEvent(lvl) }

func (l Logger) Printf(format string, v ...interface{}) {
	l.Info().Msgf(format, v...)
}

func (e *Event) Str(key, value string) *Event {
	e.fields = append(e.fields, zap.String(key, value))

	return e
}

func (e *Event) Strs(key string, values []string) *Event {
	e.fields = append(e.fields, zap.Strings(key, values))

	return e
}

func (e *Event) Interfaces(key string, vals ...interface{}) *Event {
	for _, val := range vals {
		e.fields = append(e.fields, zap.Reflect(key, val))
	}

	return e
}

func (e *Event) Err(err error) *Event {
	e.fields = append(e.fields, zap.Error(err))

	return e
}

func (e *Event) Msgf(format string, params ...interface{}) {
	msg := fmt.Sprintf(format, params...)

	l := e.Logger.Logger

	switch e.lvl {
	case zapcore.DebugLevel:
		l.Debug(msg, e.fields...)
	case zapcore.InfoLevel:
		l.Info(msg, e.fields...)
	case zapcore.WarnLevel:
		l.Warn(msg, e.fields...)
	case zapcore.ErrorLevel:
		l.Error(msg, e.fields...)
	case zapcore.DPanicLevel:
		l.DPanic(msg, e.fields...)
	case zapcore.PanicLevel:
		l.Panic(msg, e.fields...)
	case zapcore.FatalLevel:
		l.Fatal(msg, e.fields...)
	}
}

func (e Event) Msg(msg string) {
	l := e.Logger.Logger

	switch e.lvl {
	case zapcore.DebugLevel:
		l.Debug(msg, e.fields...)
	case zapcore.InfoLevel:
		l.Info(msg, e.fields...)
	case zapcore.WarnLevel:
		l.Warn(msg, e.fields...)
	case zapcore.ErrorLevel:
		l.Error(msg, e.fields...)
	case zapcore.DPanicLevel:
		l.DPanic(msg, e.fields...)
	case zapcore.PanicLevel:
		l.Panic(msg, e.fields...)
	case zapcore.FatalLevel:
		l.Fatal(msg, e.fields...)
	}
}

func (e *Event) Send() {
	l := e.Logger.WithOptions(zap.AddCallerSkip(1))
	e.Logger.Logger = l
	e.Msg("n/a message")
}
