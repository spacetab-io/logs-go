package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	jsoniter "github.com/json-iterator/go"
)

type sentryEvent struct {
	Level     string    `json:"level"`
	Msg       string    `json:"message"`
	Err       error     `json:"error"`
	Time      time.Time `json:"time"`
	Status    int       `json:"status,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Method    string    `json:"method,omitempty"`
	URL       string    `json:"url,omitempty"`
	IP        string    `json:"ip,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
	Action    string    `json:"action,omitempty"`
}

var errSkipEvent = errors.New("skip")

func sentryPush(hub *sentry.Hub, pr io.Reader) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	dec := json.NewDecoder(pr)

	for {
		var e sentryEvent

		err := dec.Decode(&e)

		if errors.Is(err, io.EOF) {
			return
		}

		if errors.Is(err, errSkipEvent) {
			continue
		}

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "unmarshaling logger failed with error %v\n", err)

			continue
		}

		hub.CaptureException(e.Err)
	}
}

func (s sentryEvent) Error() string {
	return s.Err.Error()
}
