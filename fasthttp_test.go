package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFHLogger_Printf(t *testing.T) {
	out := &bytes.Buffer{}
	_ = initLog(out)
	fhl := FHLogger{}
	fhl.Printf("some %s", "data")

	exp := "DBG |> some data <|"
	assert.Contains(t, out.String(), exp)
}
