package log_test

import (
	"bytes"
	"testing"

	log "github.com/spacetab-io/logs-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestFHLogger_Printf(t *testing.T) {
	t.Parallel()

	out := &bytes.Buffer{}
	_ = initLog(out)
	fhl := log.FHLogger{}
	fhl.Printf("some %s", "data")

	exp := "DBG |> some data <|"
	assert.Contains(t, out.String(), exp)
}
