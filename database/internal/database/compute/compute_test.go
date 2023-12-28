package compute

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

// mockgen -source=compute.go -destination=compute_mock.go -package=compute

func TestNewCompute(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	parser := NewMockparser(ctrl)
	analyzer := NewMockanalyzer(ctrl)

	compute, err := NewCompute(nil, nil, nil)
	require.Error(t, err, "query parser is invalid")
	require.Nil(t, compute)

	compute, err = NewCompute(parser, nil, nil)
	require.Error(t, err, "query analyzer is invalid")
	require.Nil(t, compute)

	compute, err = NewCompute(parser, analyzer, nil)
	require.Error(t, err, "logger is invalid")
	require.Nil(t, compute)

	compute, err = NewCompute(parser, analyzer, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, compute)
}

func TestHandleQueryWithParsingError(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	ctrl := gomock.NewController(t)
	parser := NewMockparser(ctrl)
	parser.EXPECT().
		ParseQuery(ctx, "## key").
		Return(nil, errInvalidCommand)
	analyzer := NewMockanalyzer(ctrl)

	compute, err := NewCompute(parser, analyzer, zap.NewNop())
	require.NoError(t, err)

	query, err := compute.HandleQuery(ctx, "## key")
	require.Error(t, err, errInvalidCommand)
	require.Equal(t, Query{}, query)
}

func TestHandleQueryWithAnalyzingError(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	ctrl := gomock.NewController(t)
	parser := NewMockparser(ctrl)
	parser.EXPECT().
		ParseQuery(ctx, "TRUNCATE key").
		Return([]string{"TRUNCATE", "key"}, nil)
	analyzer := NewMockanalyzer(ctrl)
	analyzer.EXPECT().
		AnalyzeQuery(ctx, []string{"TRUNCATE", "key"}).
		Return(Query{}, errInvalidCommand)

	compute, err := NewCompute(parser, analyzer, zap.NewNop())
	require.NoError(t, err)

	query, err := compute.HandleQuery(ctx, "TRUNCATE key")
	require.Error(t, err, errInvalidCommand)
	require.Equal(t, Query{}, query)
}

func TestHandleQuery(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "tx", int64(555))

	ctrl := gomock.NewController(t)
	parser := NewMockparser(ctrl)
	parser.EXPECT().
		ParseQuery(ctx, "GET key").
		Return([]string{"GET", "key"}, nil)
	analyzer := NewMockanalyzer(ctrl)
	analyzer.EXPECT().
		AnalyzeQuery(ctx, []string{"GET", "key"}).
		Return(NewQuery(GetCommandID, []string{"key"}), nil)

	compute, err := NewCompute(parser, analyzer, zap.NewNop())
	require.NoError(t, err)

	query, err := compute.HandleQuery(ctx, "GET key")
	require.NoError(t, err)
	require.Equal(t, NewQuery(GetCommandID, []string{"key"}), query)
}
