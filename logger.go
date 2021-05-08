// Package log provides a global logger for zerolog.
package log

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

const (
	defaultLevel    = "debug"
	ctxRequestIDKey = "request_id"

	defaultCallerSkipFrames    = 2
	standAloneCallerSkipFrames = 6
)

type ZLogger struct {
	zerolog.Logger
	cfg Config
}

// logger is the global logger.
var logger, _ = newZerolog(
	Config{
		Level:  defaultLevel,
		Format: FormatText,
		Caller: &CallerConfig{CallerSkipFrames: standAloneCallerSkipFrames},
	},
	os.Stdout,
)

// set global Zerolog logger.
func Init(stage string, cfg Config, serviceAlias string, serviceVersion string, w io.Writer) (err error) {
	if w == nil {
		w = os.Stdout
	}

	if cfg.Format == "" {
		cfg.Format = FormatText
	}

	if cfg.Caller == nil {
		cfg.Caller = &CallerConfig{
			CallerSkipFrames: defaultCallerSkipFrames,
		}
	}

	if cfg.Sentry == nil || !cfg.Sentry.Enable || cfg.Sentry.DSN == "" {
		logger, err = newZerolog(cfg, w)
		if err != nil {
			return fmt.Errorf("logger init newZerolog error: %w", err)
		}

		return nil
	}

	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 2 //nolint:gomnd // 2 second transport timeout

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         cfg.Sentry.DSN,
		DebugWriter: os.Stderr,
		Debug:       cfg.Sentry.Debug,
		ServerName:  serviceAlias,
		Release:     serviceVersion,
		Environment: stage,
		Transport:   sentrySyncTransport,
	})
	if err != nil {
		return fmt.Errorf("logger init raven.New error: %w", err)
	}

	h := sentry.NewHub(client, sentry.NewScope())

	pr, pw := io.Pipe()

	go sentryPush(h, pr)

	cfg.Format = FormatJSON
	logger, err = newZerolog(cfg, io.MultiWriter(w, pw))

	return err
}

func newZerolog(cfg Config, w io.Writer) (logger ZLogger, err error) {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// CallerSkipFrameCount is the number of stack frames to skip to find the caller.
	zerolog.CallerSkipFrameCount = cfg.Caller.CallerSkipFrames

	output := w

	if cfg.Format == "text" {
		// pretty print during development
		out := zerolog.ConsoleWriter{Out: w, TimeFormat: zerolog.TimeFieldFormat, NoColor: cfg.NoColor}

		out.PartsOrder = []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.MessageFieldName,
		}

		if !cfg.Caller.Disabled {
			out.PartsOrder = append(out.PartsOrder, zerolog.CallerFieldName)
		}

		out.FormatMessage = func(i interface{}) string {
			if i == nil {
				return ""
			}

			return fmt.Sprintf("|> %s <|", i)
		}

		output = out
	}

	level, err := getLevel(cfg.Level)
	if err != nil {
		return logger, err
	}

	lctx := zerolog.New(output).With().Timestamp()

	if !cfg.Caller.Disabled {
		lctx = lctx.Caller()
	}

	logger.Logger = lctx.Logger().Level(level)
	logger.cfg = cfg

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)

	return logger, nil
}

func getLevel(lvl string) (zerolog.Level, error) {
	if lvl == "" {
		lvl = defaultLevel
	}

	level, err := zerolog.ParseLevel(lvl)
	if err != nil {
		return zerolog.DebugLevel, fmt.Errorf("get level error: %w", err)
	}

	return level, nil
}

func Logger() ZLogger {
	return logger
}

// Output duplicates the global logger and sets w as its output.
func Output(w io.Writer) zerolog.Logger {
	l, _ := newZerolog(logger.cfg, w)

	return l.Logger
}

// With creates a child logger with the field added to its context.
func With() zerolog.Context {
	return logger.With()
}

// Level creates a child logger with the minimum accepted level set to level.
func Level(level zerolog.Level) zerolog.Logger {
	return logger.Level(level)
}

// LevelString creates a child logger with the minimum accepted level set to level passed as string.
func LevelString(lvl string) zerolog.Logger {
	level, _ := getLevel(lvl)

	return logger.Level(level)
}

// Sample returns a logger with the s sampler.
func Sample(s zerolog.Sampler) zerolog.Logger {
	return logger.Sample(s)
}

// Hook returns a logger with the h Hook.
func Hook(h zerolog.Hook) zerolog.Logger {
	return logger.Hook(h)
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func Err(err error) *zerolog.Event {
	return logger.Err(err)
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *zerolog.Event {
	return logger.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func WithLevel(level zerolog.Level) *zerolog.Event {
	return logger.WithLevel(level)
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *zerolog.Event {
	return logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	logger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

// Ctx returns the logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

func contextFields(ctx context.Context) (fields map[string]interface{}) {
	fields = make(map[string]interface{})
	if requestID, ok := ctx.Value(ctxRequestIDKey).(fmt.Stringer); ok && requestID.String() != "00000000-0000-0000-0000-000000000000" {
		fields[ctxRequestIDKey] = requestID.String()
	}

	if requestID, ok := ctx.Value(ctxRequestIDKey).(string); ok && requestID != "00000000-0000-0000-0000-000000000000" {
		fields[ctxRequestIDKey] = requestID
	}

	return fields
}

// With creates a child logger with the field added to its context.
func WithCtx(ctx context.Context) *zerolog.Logger {
	l := With()
	fields := contextFields(ctx)
	l2 := l.Fields(fields).Logger()

	return &l2
}
