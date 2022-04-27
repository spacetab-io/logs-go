package log_test

import (
	"bytes"
	"testing"

	log "github.com/spacetab-io/logs-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestFHLogger_Printf(t *testing.T) {
	t.Parallel()

	out := &bytes.Buffer{}
	logger, err := initLog(t, out)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	fhl := log.FHLogger{Logger: logger}
	fhl.Printf("some %s", "data")

	assert.Contains(t, out.String(), "some data")
	assert.Contains(t, out.String(), "INFO")
}
