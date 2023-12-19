package compute

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestNewAnalyzer(t *testing.T) {
	t.Parallel()

	analyzer, err := NewAnalyzer(nil)
	require.Error(t, err, "logger is invalid")
	require.Nil(t, analyzer)

	analyzer, err = NewAnalyzer(zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, analyzer)
}

func TestAnalyzeQuery(t *testing.T) {
	tests := map[string]struct {
		tokens []string
		query  Query
		err    error
	}{
		"empty tokens": {
			tokens: []string{},
			err:    errInvalidCommand,
		},
		"invalid command": {
			tokens: []string{"TRUNCATE"},
			err:    errInvalidCommand,
		},
		"invalid number arguments for set query": {
			tokens: []string{"SET", "key"},
			err:    errInvalidArguments,
		},
		"invalid number arguments for get query": {
			tokens: []string{"GET", "key", "value"},
			err:    errInvalidArguments,
		},
		"invalid number arguments for del query": {
			tokens: []string{"GET", "key", "value"},
			err:    errInvalidArguments,
		},
		"valid set query": {
			tokens: []string{"SET", "key", "value"},
			query:  NewQuery(SetCommandID, []string{"key", "value"}),
		},
		"valid get query": {
			tokens: []string{"GET", "key"},
			query:  NewQuery(GetCommandID, []string{"key"}),
		},
		"valid del query": {
			tokens: []string{"DEL", "key"},
			query:  NewQuery(DelCommandID, []string{"key"}),
		},
	}

	ctx := context.WithValue(context.Background(), "tx", int64(555))
	analyzer, err := NewAnalyzer(zap.NewNop())
	require.NoError(t, err)

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			query, err := analyzer.AnalyzeQuery(ctx, test.tokens)
			require.Equal(t, test.query, query)
			require.Equal(t, test.err, err)
		})
	}
}

func TestAnalyzeSetQuery(t *testing.T) {
	tests := map[string]struct {
		query Query
		err   error
	}{
		"empty arguments": {
			query: NewQuery(SetCommandID, []string{}),
			err:   errInvalidArguments,
		},
		"one argument": {
			query: NewQuery(SetCommandID, []string{"one"}),
			err:   errInvalidArguments,
		},
		"two arguments": {
			query: NewQuery(SetCommandID, []string{"one", "two"}),
		},
		"three arguments": {
			query: NewQuery(SetCommandID, []string{"one", "two", "three"}),
			err:   errInvalidArguments,
		},
	}

	ctx := context.WithValue(context.Background(), "tx", int64(555))
	analyzer, err := NewAnalyzer(zap.NewNop())
	require.NoError(t, err)

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := analyzer.analyzeSetQuery(ctx, test.query)
			require.Equal(t, test.err, err)
		})
	}
}

func TestAnalyzeGetQuery(t *testing.T) {
	tests := map[string]struct {
		query Query
		err   error
	}{
		"empty arguments": {
			query: NewQuery(GetCommandID, []string{}),
			err:   errInvalidArguments,
		},
		"one argument": {
			query: NewQuery(GetCommandID, []string{"one"}),
		},
		"two arguments": {
			query: NewQuery(GetCommandID, []string{"one", "two"}),
			err:   errInvalidArguments,
		},
	}

	ctx := context.WithValue(context.Background(), "tx", int64(555))
	analyzer, err := NewAnalyzer(zap.NewNop())
	require.NoError(t, err)

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := analyzer.analyzeGetQuery(ctx, test.query)
			require.Equal(t, test.err, err)
		})
	}
}

func TestAnalyzeDelQuery(t *testing.T) {
	tests := map[string]struct {
		query Query
		err   error
	}{
		"empty arguments": {
			query: NewQuery(DelCommandID, []string{}),
			err:   errInvalidArguments,
		},
		"one argument": {
			query: NewQuery(DelCommandID, []string{"one"}),
		},
		"two arguments": {
			query: NewQuery(DelCommandID, []string{"one", "two"}),
			err:   errInvalidArguments,
		},
	}

	ctx := context.WithValue(context.Background(), "tx", int64(555))
	analyzer, err := NewAnalyzer(zap.NewNop())
	require.NoError(t, err)

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := analyzer.analyzeDelQuery(ctx, test.query)
			require.Equal(t, test.err, err)
		})
	}
}
