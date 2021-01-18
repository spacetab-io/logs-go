package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNSQLogger_Output(t *testing.T) {
	out := &bytes.Buffer{}
	_ = initLog(out)

	nsqdLogger := NewNSQLogger("debug")
	err := nsqdLogger.Output(2, "some log")
	assert.NoError(t, err, "nsqdLogger.Output error")
	assert.Contains(t, out.String(), "some log")
}
