//nolint: paralleltest
package log_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/rs/zerolog"
	log "github.com/spacetab-io/logs-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("writer as buffer init", func(t *testing.T) {
		out := &bytes.Buffer{}
		err := initLog(out)
		assert.NoError(t, err, "logger init")
	})
	t.Run("nil writer init", func(t *testing.T) {
		err := initLog(nil)
		assert.NoError(t, err, "logger init")
	})
	t.Run("default format init", func(t *testing.T) {
		err := log.Init("test", log.Config{
			Level:   "trace",
			NoColor: true,
			Caller: &log.CallerConfig{
				Disabled:         false,
				CallerSkipFrames: 2,
			},
		}, "log", "v2.*.*", nil)
		assert.NoError(t, err, "logger init")
	})
	t.Run("default level init", func(t *testing.T) {
		err := log.Init("test", log.Config{
			NoColor: true,
			Caller: &log.CallerConfig{
				Disabled:         false,
				CallerSkipFrames: 2,
			},
		}, "log", "v2.*.*", nil)
		assert.NoError(t, err, "logger init")
	})
	t.Run("default caller init", func(t *testing.T) {
		err := log.Init("test", log.Config{
			NoColor: true,
			Caller:  nil,
		}, "log", "v2.*.*", nil)
		assert.NoError(t, err, "logger init")
	})
	t.Run("wrong level init", func(t *testing.T) {
		err := log.Init("test", log.Config{
			NoColor: true,
			Level:   "fart",
		}, "log", "v2.*.*", nil)
		assert.Error(t, err, "logger init")
	})
}

func TestOutput(t *testing.T) {
	err := initLog(nil)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	out := &bytes.Buffer{}
	l := log.Output(out)
	l.Trace().Msg("trace")

	exp := "TRC |> trace <|"
	assert.Contains(t, out.String(), exp)
}

func TestLevel(t *testing.T) {
	out := &bytes.Buffer{}

	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	log.Trace().Msg("test trace")

	exp := "test trace"
	assert.Contains(t, out.String(), exp)

	l2 := log.Level(zerolog.WarnLevel)
	l2.Trace().Msg("test trace")

	exp = ""
	assert.Contains(t, out.String(), exp)
}

func TestLevelString(t *testing.T) {
	out := &bytes.Buffer{}

	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	log.Trace().Msg("test trace")

	exp := "test trace"
	assert.Contains(t, out.String(), exp)

	l2 := log.LevelString("warn")
	l2.Trace().Msg("test trace")

	exp = ""
	assert.Contains(t, out.String(), exp)
}

func TestTrace(t *testing.T) {
	out := &bytes.Buffer{}

	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	log.Trace().Msg("trace")

	exp := "TRC |> trace <|"
	assert.Contains(t, out.String(), exp)
}

func TestDebug(t *testing.T) {
	out := &bytes.Buffer{}

	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	log.Debug().Msg("debug")

	exp := "DBG |> debug <|"
	assert.Contains(t, out.String(), exp)
}

func TestInfo(t *testing.T) {
	out := &bytes.Buffer{}

	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	log.Info().Msg("info")

	exp := "INF |> info <|"
	assert.Contains(t, out.String(), exp)
}

func TestWarn(t *testing.T) {
	out := &bytes.Buffer{}

	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	log.Warn().Msg("warn")

	exp := "WRN |> warn <|"
	assert.Contains(t, out.String(), exp)
}

func TestError(t *testing.T) {
	out := &bytes.Buffer{}

	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	log.Error().Msg("error")

	exp := "ERR |> error <|"
	assert.Contains(t, out.String(), exp)
}

var errForTest = errors.New("some err")

func TestErr(t *testing.T) {
	out := &bytes.Buffer{}

	err := initLog(out)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	log.Err(errForTest).Msg("error")

	exp := "ERR |> error <| "
	assert.Contains(t, out.String(), exp)

	exp = " > error=\"some err\""
	assert.Contains(t, out.String(), exp)
}

func initLog(w io.Writer) error {
	return log.Init("test", log.Config{
		Level:   "trace",
		Format:  "text",
		NoColor: true,
		Caller: &log.CallerConfig{
			Disabled:         false,
			CallerSkipFrames: 2,
		},
		Sentry: nil,
	}, "log", "v2.*.*", w)
}
