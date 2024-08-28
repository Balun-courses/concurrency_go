package compute

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewCompute(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		logger *zap.Logger

		expectedErr    error
		expectedNilObj bool
	}{
		"create compute without logger": {
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create compute": {
			logger:      zap.NewNop(),
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			compute, err := NewCompute(test.logger)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, compute)
			} else {
				assert.NotNil(t, compute)
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		queryStr string

		expectedQuery Query
		expectedErr   error
	}{
		"empty query": {
			queryStr:    "",
			expectedErr: errInvalidQuery,
		},
		"empty query without tokens": {
			queryStr:    "   ",
			expectedErr: errInvalidQuery,
		},
		"query with UTF symbols": {
			queryStr:    "字文下",
			expectedErr: errInvalidCommand,
		},
		"invalid command": {
			queryStr:    "TRUNCATE",
			expectedErr: errInvalidCommand,
		},
		"invalid number arguments for set query": {
			queryStr:    "SET key",
			expectedErr: errInvalidArguments,
		},
		"invalid number arguments for get query": {
			queryStr:    "GET key value",
			expectedErr: errInvalidArguments,
		},
		"invalid number arguments for del query": {
			queryStr:    "GET key value",
			expectedErr: errInvalidArguments,
		},
		"set query": {
			queryStr:      "SET __key__\nvalue",
			expectedQuery: NewQuery(SetCommandID, []string{"__key__", "value"}),
		},
		"get query": {
			queryStr:      "GET\t1key2",
			expectedQuery: NewQuery(GetCommandID, []string{"1key2"}),
		},
		"del query": {
			queryStr:      "DEL  /key-",
			expectedQuery: NewQuery(DelCommandID, []string{"/key-"}),
		},
	}

	compute, err := NewCompute(zap.NewNop())
	require.NoError(t, err)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			query, err := compute.Parse(test.queryStr)
			assert.Equal(t, test.expectedErr, err)
			assert.True(t, reflect.DeepEqual(test.expectedQuery, query))
		})
	}
}
