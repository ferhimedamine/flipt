package server

import (
	"context"
	"errors"
	"testing"

	"github.com/markphelps/flipt/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorInterceptor(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode codes.Code
	}{
		{
			name:     "not found error",
			err:      storage.ErrNotFound("foo"),
			wantCode: codes.NotFound,
		},
		{
			name:     "invalid error",
			err:      storage.ErrInvalid("foo"),
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "other error",
			err:      errors.New("foo"),
			wantCode: codes.Internal,
		},
		{
			name: "no error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			spy := grpc.UnaryHandler(func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, tt.err
			})

			_, err := ErrorInterceptor(context.Background(), nil, nil, spy)
			if tt.err != nil {
				require.Error(t, err)
				status := status.Convert(err)
				assert.Equal(t, tt.wantCode, status.Code())
				return
			}

			require.NoError(t, err)
		})
	}
}
