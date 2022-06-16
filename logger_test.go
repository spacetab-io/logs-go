package log_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	cfgstructs "github.com/spacetab-io/configuration-structs-go/v2"
	"github.com/spacetab-io/configuration-structs-go/v2/contracts"
	log "github.com/spacetab-io/logs-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestInitNew(t *testing.T) {
	t.Parallel()

	l, err := log.Init(&cfgstructs.Logs{
		Level:   "debug",
		Format:  "text",
		Colored: true,
		Caller:  cfgstructs.CallerConfig{Show: true, SkipFrames: 1},
		Sentry: &cfgstructs.SentryConfig{
			Enable: true,
			Debug:  true,
			DSN:    os.Getenv("SENTRY_DSN"),
		},
	}, "test", "service", "v1.0.0", os.Stdout)
	if !assert.NoError(t, err, "logger init") {
		t.FailNow()
	}

	l.Debug().Str("some", "string").Str("and", "more").Msgf("%s!", "hello")
	l.Info().Strs("array", []string{"string", "and", "more"}).Msg("another!")
	l.Warn().Interfaces("another", "string", 1, 1.2, struct{ Boo string }{Boo: "boo"}).Str("with", "more").Msgf("%s", "world!")
	l.Error().Err(log.ErrEmptyOutput).Str("another", "string").Str("with", "more").Msgf("%s %s", "catch", "it!")
	l.Error().Err(log.ErrEmptyOutput).Send()
}

func TestInit(t *testing.T) {
	type inStruct struct {
		stage          string
		cfg            contracts.LogsCfgInterface
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
