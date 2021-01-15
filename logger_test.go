package log

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrace(t *testing.T) {
	out := &bytes.Buffer{}
	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	Trace().Msg("trace")
	exp := "TRC |> trace <|"
	assert.Contains(t, out.String(), exp)
}

func TestDebug(t *testing.T) {
	out := &bytes.Buffer{}
	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	Debug().Msg("debug")
	exp := "DBG |> debug <|"
	assert.Contains(t, out.String(), exp)
}

func TestInfo(t *testing.T) {
	out := &bytes.Buffer{}
	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	Info().Msg("info")
	exp := "INF |> info <|"
	assert.Contains(t, out.String(), exp)
}

func TestWarn(t *testing.T) {
	out := &bytes.Buffer{}
	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	Warn().Msg("warn")
	exp := "WRN |> warn <|"
	assert.Contains(t, out.String(), exp)
}

func TestError(t *testing.T) {
	out := &bytes.Buffer{}
	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	Error().Msg("error")
	exp := "ERR |> error <|"
	assert.Contains(t, out.String(), exp)
}

func TestErr(t *testing.T) {
	out := &bytes.Buffer{}
	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	Error().Err(errors.New("some err")).Msg("error")
	exp := "ERR |> error <| logger_test.go:79 > error=\"some err\""
	assert.Contains(t, out.String(), exp)
}

func initLog(w io.Writer) error {
	return Init("test", Config{
		Level:      "trace",
		Format:     "text",
		NoColor:    true,
		ShowCaller: true,
		Sentry:     nil,
	}, "log", "v2.*.*", w)
}
