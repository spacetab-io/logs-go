//nolint: paralleltest
package log_test

import (
	"bytes"
	"testing"

	log "github.com/spacetab-io/logs-go/v2"
	"github.com/stretchr/testify/assert"
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
			name: "uncommon err",
			in:   "some log",
			exp:  "some log",
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			_ = initLog(out)

			nsqdLogger := log.NewNSQLogger("info")

			err := nsqdLogger.Output(2, tc.in)
			if !assert.NoError(t, err, "nsqdLogger.Output error") {
				t.FailNow()
			}

			assert.Contains(t, out.String(), tc.exp)
		})
	}
}
