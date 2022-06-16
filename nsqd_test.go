//nolint: paralleltest
package log_test

import (
	"bytes"
	"testing"

	log "github.com/spacetab-io/logs-go/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNSQLogger_Output(t *testing.T) {
	type tc struct {
		name string
		in   string
		exp  string
	}

	tcs := []tc{
		{
			name: "NSQ error log",
			in:   "ERR    1 (localhost:4150) error connecting to nsqd - dial tcp [::1]:4150: connect: connection refused",
			exp:  "1 (localhost:4150) error connecting to nsqd - dial tcp [::1]:4150: connect: connection refused",
		},
		{
			name: "NSQ trace log",
			in:   "TRC some log",
			exp:  "some log",
		},
		{
			name: "NSQ info log",
			in:   "INF some log",
			exp:  "some log",
		},
		{
			name: "NSQ warning log",
			in:   "WRN some log",
			exp:  "some log",
		},
		//{
		//	name: "NSQ fatal log",
		//	in:   "FTL some log",
		//	exp:  "some log",
		// },
		//{
		//	name: "NSQ panic log",
		//	in:   "PNC some log",
		//	exp:  "some log",
		//},
		{
			name: "uncommon err",
			in:   "some log",
			exp:  "some log",
		},
		{
			name: "not an error",
			in:   "",
			exp:  "",
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			logger, _ := initLog(t, out)

			nsqdLogger := log.NewNSQLogger(logger)

			defer func() {
				if err := recover(); err != nil {
					assert.Contains(t, out.String(), tc.exp)
				}
			}()

			err := nsqdLogger.Output(2, tc.in)
			if !assert.NoError(t, err, "nsqdLogger.Output error") {
				t.FailNow()
			}

			assert.Contains(t, out.String(), tc.exp)
		})
	}
}

func TestNSQLogger_LogLevel(t *testing.T) {
	type tc struct {
		name string
		in   zapcore.Level
		exp  int
	}

	tcs := []tc{
		{
			name: zapcore.DebugLevel.String(),
			in:   zapcore.DebugLevel,
			exp:  0,
		},
		{
			name: zapcore.InfoLevel.String(),
			in:   zapcore.InfoLevel,
			exp:  1,
		},
		{
			name: zapcore.WarnLevel.String(),
			in:   zapcore.WarnLevel,
			exp:  2,
		},
		{
			name: zapcore.ErrorLevel.String(),
			in:   zapcore.ErrorLevel,
			exp:  3,
		},
		{
			name: zapcore.FatalLevel.String(),
			in:   zapcore.FatalLevel,
			exp:  3,
		},
		{
			name: zapcore.PanicLevel.String(),
			in:   zapcore.PanicLevel,
			exp:  3,
		},
		{
			name: zapcore.DPanicLevel.String(),
			in:   zapcore.DPanicLevel,
			exp:  3,
		},
		{
			name: zapcore.Level(-2).String(),
			in:   zapcore.Level(-2),
			exp:  4,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := log.NewNSQLogger(log.Logger{Level: tc.in})
			assert.Equal(t, tc.exp, l.LogLevel())
		})
	}
}
