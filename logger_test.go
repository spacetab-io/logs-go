package log_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	cfgstructs "github.com/spacetab-io/configuration-structs-go"
	log "github.com/spacetab-io/logs-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestInitNew(t *testing.T) {

	l, err := log.Init(&cfgstructs.Logs{
		Level:   "debug",
		Format:  "text",
		Colored: true,
		Caller:  cfgstructs.CallerConfig{Show: true, SkipFrames: 1},
		Sentry: &cfgstructs.SentryConfig{
			Enable: true,
			Debug:  true,
			DSN:    "https://60acad048aee4777a7331d8282fe8022@sentry.spacetab.io/2",
		},
	}, "test", "service", "v1.0.0", os.Stdout)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	l.Info().Str("some", "string").Str("and", "more").Msg("hello!")
	l.Warn().Str("another", "string").Str("with", "more").Msg("world!")
	l.Error().Err(log.ErrEmptyOutput).Str("another", "string").Str("with", "more").Msg("catch it!")
	l.Error().Err(log.ErrEmptyOutput).Send()

	//assert.Contains(t, out.String(), "string")
	//assert.Contains(t, out.String(), "more")
}

func TestInit(t *testing.T) {
	type inStruct struct {
		stage          string
		cfg            cfgstructs.LogsInterface
		serviceAlias   string
		serviceVersion string
		w              bool
		msg            string
	}

	type testCase struct {
		name   string
		in     inStruct
		hasErr bool
		exp    string
	}

	testCfg := getTestCfg(t)
	tcs := []testCase{
		{
			name: "writer as buffer init",
			in: inStruct{
				stage:          "test",
				cfg:            testCfg,
				serviceAlias:   "log",
				serviceVersion: "v0.0.0",
				w:              true,
				msg:            "some warn",
			},
			exp: "some warn",
		},
		{
			name: "nil writer init",
			in: inStruct{
				stage:          "test",
				cfg:            testCfg,
				serviceAlias:   "log",
				serviceVersion: "v0.0.0",
				w:              false,
				msg:            "some warn",
			},
			hasErr: true,
			// err: log.ErrEmptyOutput,
		},
		{
			name: "default format init",
			in: inStruct{
				stage: "test",
				cfg: &cfgstructs.Logs{
					Level:   "debug",
					Colored: true,
					Caller: cfgstructs.CallerConfig{
						Show:       false,
						SkipFrames: 2,
					},
				},
				serviceAlias:   "log",
				serviceVersion: "v0.0.0",
				w:              true,
				msg:            "some warn",
			},
			exp: "some warn",
		},
		{
			name: "default level init",
			in: inStruct{
				stage: "test",
				cfg: &cfgstructs.Logs{
					Colored: true,
					Caller: cfgstructs.CallerConfig{
						Show:       false,
						SkipFrames: 2,
					},
				},
				serviceAlias:   "log",
				serviceVersion: "v0.0.0",
				w:              true,
				msg:            "some warn",
			},
			exp: "some warn",
		},
		{
			name: "default caller init",
			in: inStruct{
				stage: "test",
				cfg: &cfgstructs.Logs{
					Colored: true,
					Caller:  cfgstructs.CallerConfig{},
				},
				serviceAlias:   "log",
				serviceVersion: "v0.0.0",
				w:              true,
				msg:            "some warn",
			},
			exp: "some warn",
		},
		{
			name: "wrong level init",
			in: inStruct{
				stage: "test",
				cfg: &cfgstructs.Logs{
					Colored: true,
					Level:   "fart",
				},
				serviceAlias:   "log",
				serviceVersion: "v0.0.0",
				w:              true,
			},
			hasErr: true,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var w io.Writer

			if tc.in.w {
				w = &bytes.Buffer{}
			}

			l, err := log.Init(tc.in.cfg, tc.in.stage, tc.in.serviceAlias, tc.in.serviceVersion, w)
			if tc.hasErr {
				if !assert.Error(t, err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}

				l.Debug().Msg(tc.in.msg)
				assert.Contains(t, w.(*bytes.Buffer).String(), tc.exp)
				assert.Contains(t, w.(*bytes.Buffer).String(), "DEBUG")

				l.Info().Msg(tc.in.msg)
				assert.Contains(t, w.(*bytes.Buffer).String(), "INFO")

				l.Warn().Msg(tc.in.msg)
				assert.Contains(t, w.(*bytes.Buffer).String(), "WARN")

				errForTest := errors.New("some err")

				l.Error().Err(errForTest).Msg(tc.in.msg)
				assert.Contains(t, w.(*bytes.Buffer).String(), "ERROR")
				assert.Contains(t, w.(*bytes.Buffer).String(), "\"error\": \"some err\"")

				l.DPanic().Msg(tc.in.msg)
				assert.Contains(t, w.(*bytes.Buffer).String(), "DPANIC")

				//assert.Contains(t, w.(*bytes.Buffer).String(), "PANIC")
			}
		})
	}
}

func initLog(t *testing.T, w io.Writer) (log.Logger, error) {
	t.Helper()

	return log.Init(getTestCfg(t), "test", "log", "v2.*.*", w)
}

func getTestCfg(t *testing.T) *cfgstructs.Logs {
	t.Helper()

	return &cfgstructs.Logs{
		Level:   "debug",
		Format:  "text",
		Colored: false,
		Caller: cfgstructs.CallerConfig{
			Show:       true,
			SkipFrames: 2,
		},
		Sentry: nil,
	}
}
