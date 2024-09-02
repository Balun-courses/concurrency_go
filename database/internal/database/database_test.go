package database

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"spider/internal/database/compute"
	"spider/internal/database/storage"
)

// mockgen -source=database.go -destination=database_mock.go -package=database

func TestNewDatabase(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		compute computeLayer
		storage storageLayer
		logger  *zap.Logger

		expectedErr    error
		expectedNilObj bool
	}{
		"create database without compute layer": {
			expectedErr:    errors.New("compute is invalid"),
			expectedNilObj: true,
		},
		"create database without storage layer": {
			compute:        NewMockcomputeLayer(ctrl),
			expectedErr:    errors.New("storage is invalid"),
			expectedNilObj: true,
		},
		"create database without logger": {
			compute:        NewMockcomputeLayer(ctrl),
			storage:        NewMockstorageLayer(ctrl),
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create database": {
			compute: NewMockcomputeLayer(ctrl),
			storage: NewMockstorageLayer(ctrl),
			logger:  zap.NewNop(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			storage, err := NewDatabase(test.compute, test.storage, test.logger)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, storage)
			} else {
				assert.NotNil(t, storage)
			}
		})
	}
}

func TestHandleQuery(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		query        string
		computeLayer func() computeLayer
		storageLayer func() storageLayer

		expectedResponse string
	}{
		"handle incorrect query": {
			query: "TRUNCATE",
			computeLayer: func() computeLayer {
				computeLayer := NewMockcomputeLayer(ctrl)
				computeLayer.EXPECT().
					Parse("TRUNCATE").
					Return(compute.Query{}, errors.New("compute error"))
				return computeLayer
			},
			storageLayer:     func() storageLayer { return NewMockstorageLayer(ctrl) },
			expectedResponse: "[error] compute error",
		},
		"handle set query with error from storage": {
			query: "SET key value",
			computeLayer: func() computeLayer {
				computeLayer := NewMockcomputeLayer(ctrl)
				computeLayer.EXPECT().
					Parse("SET key value").
					Return(compute.NewQuery(
						compute.SetCommandID,
						[]string{"key", "value"},
					), nil)
				return computeLayer
			},
			storageLayer: func() storageLayer {
				storageLayer := NewMockstorageLayer(ctrl)
				storageLayer.EXPECT().
					Set(gomock.Any(), "key", "value").
					Return(errors.New("storage error"))
				return storageLayer
			},
			expectedResponse: "[error] storage error",
		},
		"handle set query": {
			query: "SET key value",
			computeLayer: func() computeLayer {
				computeLayer := NewMockcomputeLayer(ctrl)
				computeLayer.EXPECT().
					Parse("SET key value").
					Return(compute.NewQuery(
						compute.SetCommandID,
						[]string{"key", "value"},
					), nil)
				return computeLayer
			},
			storageLayer: func() storageLayer {
				storageLayer := NewMockstorageLayer(ctrl)
				storageLayer.EXPECT().
					Set(gomock.Any(), "key", "value").
					Return(nil)
				return storageLayer
			},
			expectedResponse: "[ok]",
		},
		"handle del query with error from storage": {
			query: "DEL key",
			computeLayer: func() computeLayer {
				computeLayer := NewMockcomputeLayer(ctrl)
				computeLayer.EXPECT().
					Parse("DEL key").
					Return(compute.NewQuery(
						compute.DelCommandID,
						[]string{"key"},
					), nil)
				return computeLayer
			},
			storageLayer: func() storageLayer {
				storageLayer := NewMockstorageLayer(ctrl)
				storageLayer.EXPECT().
					Del(gomock.Any(), "key").
					Return(errors.New("storage error"))
				return storageLayer
			},
			expectedResponse: "[error] storage error",
		},
		"handle del query": {
			query: "DEL key",
			computeLayer: func() computeLayer {
				computeLayer := NewMockcomputeLayer(ctrl)
				computeLayer.EXPECT().
					Parse("DEL key").
					Return(compute.NewQuery(
						compute.DelCommandID,
						[]string{"key"},
					), nil)
				return computeLayer
			},
			storageLayer: func() storageLayer {
				storageLayer := NewMockstorageLayer(ctrl)
				storageLayer.EXPECT().
					Del(gomock.Any(), "key").
					Return(nil)
				return storageLayer
			},
			expectedResponse: "[ok]",
		},
		"handle get query with error from storage": {
			query: "GET key",
			computeLayer: func() computeLayer {
				computeLayer := NewMockcomputeLayer(ctrl)
				computeLayer.EXPECT().
					Parse("GET key").
					Return(compute.NewQuery(
						compute.GetCommandID,
						[]string{"key"},
					), nil)
				return computeLayer
			},
			storageLayer: func() storageLayer {
				storageLayer := NewMockstorageLayer(ctrl)
				storageLayer.EXPECT().
					Get(gomock.Any(), "key").
					Return("", errors.New("storage error"))
				return storageLayer
			},
			expectedResponse: "[error] storage error",
		},
		"handle get query with not found error from storage": {
			query: "GET key",
			computeLayer: func() computeLayer {
				computeLayer := NewMockcomputeLayer(ctrl)
				computeLayer.EXPECT().
					Parse("GET key").
					Return(compute.NewQuery(
						compute.GetCommandID,
						[]string{"key"},
					), nil)
				return computeLayer
			},
			storageLayer: func() storageLayer {
				storageLayer := NewMockstorageLayer(ctrl)
				storageLayer.EXPECT().
					Get(gomock.Any(), "key").
					Return("", storage.ErrorNotFound)
				return storageLayer
			},
			expectedResponse: "[not found]",
		},
		"handle get query": {
			query: "GET key",
			computeLayer: func() computeLayer {
				computeLayer := NewMockcomputeLayer(ctrl)
				computeLayer.EXPECT().
					Parse("GET key").
					Return(compute.NewQuery(
						compute.GetCommandID,
						[]string{"key"},
					), nil)
				return computeLayer
			},
			storageLayer: func() storageLayer {
				storageLayer := NewMockstorageLayer(ctrl)
				storageLayer.EXPECT().
					Get(gomock.Any(), "key").
					Return("value", nil)
				return storageLayer
			},
			expectedResponse: "[ok] value",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			storage, err := NewDatabase(test.computeLayer(), test.storageLayer(), zap.NewNop())
			require.NoError(t, err)

			response := storage.HandleQuery(context.Background(), test.query)
			assert.Equal(t, test.expectedResponse, response)
		})
	}
}
