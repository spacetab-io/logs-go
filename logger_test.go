package logs

import (
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	testFormatter = &logrus.TextFormatter{
		TimestampFormat:        time.RFC3339,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		QuoteEmptyFields:       true,
	}
)

// I just want to test it before pushing. Don't know how to test it in right way, sorry )
func TestNewLogger(t *testing.T) {
	type testCase struct {
		name        string
		config      *Config
		expFormatte *logrus.TextFormatter
		hasErr      bool
	}
	type testCases []testCase

	cfg := &Config{
		Stage:    "test",
		LogLevel: logrus.InfoLevel.String(),
		Sentry: &SentryConfig{
			Enable: false,
		},
	}
	cfg2 := &Config{
		Stage:    "test",
		LogLevel: logrus.InfoLevel.String(),
		Sentry: &SentryConfig{
			Enable: true,
			DSN:    "https://xxx@sentry.io/yyy",
		},
	}

	cfg3 := &Config{
		Stage:    "",
		LogLevel: logrus.InfoLevel.String(),
		Sentry: &SentryConfig{
			Enable: true,
			DSN:    "go away",
		},
	}
	cfg4 := &Config{
		Stage:    "test",
		LogLevel: logrus.InfoLevel.String(),
		Sentry: &SentryConfig{
			Enable: true,
			DSN:    "go away",
		},
	}
	cfg5 := &Config{
		Stage:    "test",
		LogLevel: "paranoya",
		Sentry: &SentryConfig{
			Enable: true,
			DSN:    "go away",
		},
	}

	tcs := testCases{
		{name: "new logger", config: cfg, expFormatte: testFormatter},
		{name: "good sentry", config: cfg2, expFormatte: testFormatter},
		{name: "no stage", config: cfg3, hasErr: true},
		{name: "bad sentry dsn", config: cfg4, hasErr: true},
		{name: "bad sentry log lvl", config: cfg5, hasErr: true},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			l, err := NewLogger(tc.config)

			if tc.hasErr {
				if !assert.Error(t, err, "has error fail") {
					t.FailNow()
				}
				return
			} else {
				if !assert.NoError(t, err, "no error fail") {
					t.FailNow()
				}
			}

			if !assert.Equal(t, tc.config.LogLevel, l.GetLevel().String(), "wrong log level") {
				t.FailNow()
			}

			if !assert.Equal(t, tc.expFormatte, l.Formatter, "differ formatter") {
				t.FailNow()
			}

			if tc.config.Sentry.Enable {
				if !assert.Equal(t, 4, len(l.Hooks), "wrong number of hooks") {
					t.FailNow()
				}
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	t.Run("get ent key value", func(t *testing.T) {
		c := &Config{}
		key := "STAGE"
		val := "testing"
		err := os.Setenv(key, val)
		assert.NoError(t, err)
		c.SetStage()
		assert.Equal(t, val, c.Stage)
		_ = os.Unsetenv(key)
	})
	t.Run("get fallback stage", func(t *testing.T) {
		c := &Config{}
		key := "notSTAGE"
		val := "testing"
		err := os.Setenv(key, val)
		assert.NoError(t, err)
		c.SetStage()
		assert.Equal(t, "development", c.Stage)
	})
}

func TestConfig_SetStage(t *testing.T) {
	type testCase struct {
		name   string
		key    string
		flb    string
		val    string
		hasErr bool
	}
	type testCases []testCase

	tcs := testCases{
		{name: "get env value", key: "key", flb: "", val: "value", hasErr: false},
		{name: "get fallback value", key: "", flb: "fallback", val: "fallback", hasErr: false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(tc.key, tc.val)
			val := GetEnv(tc.key, tc.flb)
			if !assert.Equal(t, tc.val, val) {
				t.FailNow()
			}
		})
	}
}
