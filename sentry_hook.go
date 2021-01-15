package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/getsentry/raven-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

type errWithStackTrace struct {
	Err        string            `json:"error"`
	Stacktrace *raven.Stacktrace `json:"stacktrace"`
}

type sentryEvent struct {
	Level     string            `json:"level"`
	Msg       string            `json:"message"`
	Err       errWithStackTrace `json:"error"`
	Time      time.Time         `json:"time"`
	Status    int               `json:"status,omitempty"`
	UserAgent string            `json:"user_agent,omitempty"`
	Method    string            `json:"method,omitempty"`
	URL       string            `json:"url,omitempty"`
	IP        string            `json:"ip,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	Action    string            `json:"action,omitempty"`
}

var errSkipEvent = errors.New("skip")

func sentryPush(stage string, serviceAlias string, serviceVersion string, client *raven.Client, pr io.Reader) {
	defer client.Close()

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	dec := json.NewDecoder(pr)

	for {
		var e sentryEvent

		err := dec.Decode(&e)

		switch err {
		case nil:
			break
		case io.EOF:
			return
		case errSkipEvent:
			continue
		default:
			_, _ = fmt.Fprintf(os.Stderr, "unmarshaling logger failed with error %v\n", err)
			continue
		}

		var level raven.Severity

		switch e.Level {
		case "debug":
			level = raven.DEBUG
		case "info":
			level = raven.INFO
		case "warn":
			level = raven.WARNING
		case "error":
			level = raven.ERROR
		case "fatal", "panic":
			level = raven.FATAL
		default:
			continue
		}

		packet := raven.Packet{
			Message:     e.Msg,
			Timestamp:   raven.Timestamp(e.Time),
			Level:       level,
			Platform:    "go",
			Project:     serviceAlias,
			Logger:      "zerolog",
			Release:     serviceVersion,
			Culprit:     e.Err.Err,
			Environment: stage,
		}

		if e.Err.Stacktrace != nil {
			packet.Interfaces = append(packet.Interfaces, e.Err.Stacktrace)
		}

		if e.IP != "" {
			packet.Interfaces = append(packet.Interfaces, &raven.User{IP: e.IP})
		}

		if e.URL != "" {
			h := raven.Http{
				URL:     e.URL,
				Method:  e.Method,
				Headers: make(map[string]string),
			}
			if e.UserAgent != "" {
				h.Headers["User-Agent"] = e.UserAgent
			}

			packet.Interfaces = append(packet.Interfaces, &h)
		}

		_, _ = client.Capture(&packet, nil)
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func stackTraceToSentry(st errors.StackTrace) *raven.Stacktrace {
	var frames []*raven.StacktraceFrame

	for _, f := range st {
		pc := uintptr(f) - 1
		fn := runtime.FuncForPC(pc)

		var (
			funcName, file string
			line           int
		)

		unk := "unknown"

		if fn != nil {
			file, line = fn.FileLine(pc)
			funcName = fn.Name()
		} else {
			file = unk
			funcName = unk
		}

		frame := raven.NewStacktraceFrame(pc, funcName, file, line, 3, nil)
		if frame != nil {
			frames = append([]*raven.StacktraceFrame{frame}, frames...)
		}
	}

	return &raven.Stacktrace{Frames: frames}
}
