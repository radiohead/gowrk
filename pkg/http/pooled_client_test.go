package http_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	wrkHttp "github.com/radiohead/gowrk/pkg/http"
	"github.com/radiohead/gowrk/test/mocks"
)

func TestNewPooledClient(t *testing.T) {
	tests := []struct {
		name     string
		poolSize uint32
		wantErr  bool
	}{
		{
			name:    "when pool size is 0",
			wantErr: true,
		},
		{
			name:     "when pool size is 100",
			poolSize: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := wrkHttp.NewPooledClient(tt.poolSize, 10)

			if tt.wantErr {
				assert.Nil(t, client)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, client)
				assert.NoError(t, err)
			}
		})
	}
}

func TestPooledClientDo(t *testing.T) {
	tests := []struct {
		name    string
		request *http.Request
		want    *http.Response
		wantErr error
	}{
		{
			name:    "when underlying client returns an error",
			request: &http.Request{},
			wantErr: errors.New("http client error"),
		},
		{
			name:    "when underlying client returns response",
			request: &http.Request{},
			want:    &http.Response{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			clientPool := make([]wrkHttp.Client, 2)
			for i := 0; i < 2; i++ {
				mock := mocks.NewMockClient(ctrl)
				mock.EXPECT().Do(tt.request).Return(tt.want, tt.wantErr)

				clientPool[i] = mock
			}

			pooledClient, err := wrkHttp.NewPooledClientWithPool(clientPool)
			assert.NoError(t, err)

			for i := 0; i < 2; i++ {
				res, err := pooledClient.Do(tt.request)
				assert.Equal(t, tt.want, res)
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}
