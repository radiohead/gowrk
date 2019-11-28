package http_test

import (
	"net/http"
	"testing"

	wrkHttp "github.com/radiohead/gowrk/pkg/http"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		connCount int
	}{
		{
			name:      "when connection count is 10",
			connCount: 10,
		},
		{
			name:      "when connection count is 100",
			connCount: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := wrkHttp.NewClient(uint32(tt.connCount))

			tr, ok := c.Transport.(*http.Transport)
			if !ok {
				assert.FailNow(t, "client transport is not http.Transport", tr)
			}

			assert.Equal(t, tt.connCount, tr.MaxIdleConns)
			assert.Equal(t, tt.connCount, tr.MaxIdleConnsPerHost)
		})
	}
}
