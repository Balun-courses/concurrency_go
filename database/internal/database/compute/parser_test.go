package compute

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestNewParser(t *testing.T) {
	t.Parallel()

	parser, err := NewParser(nil)
	require.Error(t, err, "logger is invalid")
	require.Nil(t, parser)

	parser, err = NewParser(zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, parser)
}

func TestParse(t *testing.T) {
	tests := map[string]struct {
		query  string
		tokens []string
		err    error
	}{
		"empty query": {
			query: "",
		},
		"query without tokens": {
			query: "   ",
		},
		"query with UTF symbols": {
			query: "字文下",
			err:   errInvalidSymbol,
		},
		"query with one token": {
			query:  "set",
			tokens: []string{"set"},
		},
		"query with two tokens": {
			query:  "set key",
			tokens: []string{"set", "key"},
		},
		"query with one token with digits": {
			query:  "2set1",
			tokens: []string{"2set1"},
		},
		"query with one token with underscores": {
			query:  "_set__",
			tokens: []string{"_set__"},
		},
		"query with one token with invalid symbols": {
			query: ".set#",
			err:   errInvalidSymbol,
		},
		"query with two tokens with additional spaces": {
			query:  " set   key  ",
			tokens: []string{"set", "key"},
		},
	}

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			parser, err := NewParser(zap.NewNop())
			require.NoError(t, err)

			tokens, err := parser.ParseQuery(ctx, test.query)
			require.Equal(t, test.err, err)
			require.True(t, reflect.DeepEqual(test.tokens, tokens))
		})
	}
}
