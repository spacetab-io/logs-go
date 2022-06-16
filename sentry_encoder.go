package log

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type sentryEncoder struct {
	zapcore.Encoder
	client *sentry.Client
}

// severityMap is a mapping of logrus log level to sentry log level.
var severityMap = map[zapcore.Level]sentry.Level{
	zapcore.DebugLevel: sentry.LevelDebug,
	zapcore.InfoLevel:  sentry.LevelInfo,
	zapcore.WarnLevel:  sentry.LevelWarning,
	zapcore.ErrorLevel: sentry.LevelError,
	zapcore.FatalLevel: sentry.LevelFatal,
	zapcore.PanicLevel: sentry.LevelFatal,
}

// SentryEventIdentityModifier is a sentry event modifier that simply passes
// through the event.
type SentryEventIdentityModifier struct{}

// ApplyToEvent simply returns the event (ignoring the hint).
func (m *SentryEventIdentityModifier) ApplyToEvent(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	return event
}

var sentryModifier = &SentryEventIdentityModifier{}

func (e *sentryEncoder) Clone() zapcore.Encoder {
	return e.Encoder.Clone()
}

func (e *sentryEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	if entry.Level < zapcore.ErrorLevel {
		buf, err := e.Encoder.EncodeEntry(entry, fields)
		if err != nil {
			return nil, fmt.Errorf("encode entry error: %w", err)
		}

		return buf, nil
	}

	event := sentry.NewEvent()
	event.Message = entry.Message
	event.Level = severityMap[entry.Level]
	event.Timestamp = entry.Time

	final := zapcore.NewMapObjectEncoder()

	for _, field := range fields {
		field.AddTo(final)
	}

	for k, v := range final.Fields {
		if k == "error" {
			st := sentry.NewStacktrace()
			st.Frames = append(st.Frames, frameFromCaller(entry.Caller))
			event.Exception = []sentry.Exception{{
				Type:       entry.Message,
				Value:      entry.Caller.String(),
				Stacktrace: st,
			}}
		} else {
			event.Extra[k] = v
		}
	}

	e.client.CaptureEvent(event, nil, sentryModifier)

	e.client.Flush(1 * time.Second)

	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, fmt.Errorf("encode entry error: %w", err)
	}

	return buf, nil
}

func frameFromCaller(caller zapcore.EntryCaller) sentry.Frame {
	return sentry.Frame{
		Function: caller.Function,
		Module:   caller.Function,
		Filename: caller.File,
		AbsPath:  caller.FullPath(),
		Package:  caller.TrimmedPath(),
		Lineno:   caller.Line,
		InApp:    caller.Defined,
	}
}
