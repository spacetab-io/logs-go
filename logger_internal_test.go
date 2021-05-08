package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_contextFields(t *testing.T) {
	type tc struct {
		name string
		in   context.Context
		exp  map[string]interface{}
	}

	rID := newTestUUID("45bf025d-9e46-45a4-8562-c37c4d48a9ca")

	tcs := []tc{
		{
			name: "context with requers_id",
			in:   context.WithValue(context.Background(), ctxRequestIDKey, rID), //nolint:staticcheck // да ладно!
			exp:  map[string]interface{}{ctxRequestIDKey: rID.String()},
		},
		{
			name: "empty context",
			in:   context.Background(),
			exp:  map[string]interface{}{},
		},
		{
			name: "string equest id",
			in:   context.WithValue(context.Background(), ctxRequestIDKey, rID.String()), //nolint:staticcheck // ну хорош!
			exp:  map[string]interface{}{ctxRequestIDKey: rID.String()},
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := contextFields(tc.in)

			assert.Equal(t, tc.exp, out)
		})
	}
}

type testUUID struct {
	str string
}

func (t testUUID) String() string {
	return t.str
}

func newTestUUID(str string) testUUID {
	return testUUID{str: str}
}
