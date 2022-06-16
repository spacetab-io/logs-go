package log

import (
	"context"
	"testing"

	"github.com/spacetab-io/configuration-structs-go/v2/contracts"
	"github.com/stretchr/testify/assert"
)

func Test_contextFields(t *testing.T) {
	type tc struct {
		name string
		in   context.Context
		exp  map[contracts.ContextKey]interface{}
	}

	tcs := []tc{
		{
			name: "context with requers_id",
			in:   context.WithValue(context.TODO(), contracts.ContextKeyRequestID, newTestUUID("45bf025d-9e46-45a4-8562-c37c4d48a9cb")),
			exp:  map[contracts.ContextKey]interface{}{contracts.ContextKeyRequestID: newTestUUID("45bf025d-9e46-45a4-8562-c37c4d48a9cb").String()},
		},
		{
			name: "empty context",
			in:   context.Background(),
			exp:  map[contracts.ContextKey]interface{}{},
		},
		{
			name: "string request id",
			in:   context.WithValue(context.TODO(), contracts.ContextKeyRequestID, newTestUUID("45bf025d-9e46-45a4-8562-c37c4d48a9cc").String()),
			exp:  map[contracts.ContextKey]interface{}{contracts.ContextKeyRequestID: newTestUUID("45bf025d-9e46-45a4-8562-c37c4d48a9cc").String()},
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
