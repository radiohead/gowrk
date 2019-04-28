package http

import (
	"errors"
	"math/rand"
	"net/http"

	"github.com/radiohead/gowrk/pkg/types"
)

// ErrEmptyPool is returned if the client pool is zero-size.
var ErrEmptyPool = errors.New("pool of size 0 is not allowed")

// PooledClient contains a pool of http.Client instances.
// The client implements a reduced interface of http.Client.
type PooledClient struct {
	poolSize   int
	clientPool []types.HTTPClient
}

// NewPooledClient returns a new instance of PooledClient, ready to be used.
// poolSize controls the amount of clients in the pool
// and maxIdleConns controls the size of the connection pool of individual clients.
func NewPooledClient(poolSize uint16, maxIdleConns uint16) (*PooledClient, error) {
	clients := make([]types.HTTPClient, poolSize)
	for i := 0; i < int(poolSize); i++ {
		clients[i] = NewClient(maxIdleConns)
	}

	return NewPooledClientWithPool(clients)
}

// NewPooledClientWithPool returns a new instance of PooledClient,
// with client pool set from the pool, which is copied and can be later re-used safely.
func NewPooledClientWithPool(pool []types.HTTPClient) (*PooledClient, error) {
	poolSize := len(pool)
	if poolSize < 1 {
		return nil, ErrEmptyPool
	}

	clientPool := make([]types.HTTPClient, poolSize)
	copy(clientPool, pool)

	return &PooledClient{
		clientPool: clientPool,
		poolSize:   poolSize,
	}, nil
}

// Do sends the request using a randomly selected Client from the pool.
// If the pool is empty, an error is returned.
func (c *PooledClient) Do(request *http.Request) (*http.Response, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	return client.Do(request)
}

func (c *PooledClient) getClient() (types.HTTPClient, error) {
	if c.poolSize < 1 {
		return nil, ErrEmptyPool
	}

	if c.poolSize == 1 {
		return c.clientPool[0], nil
	}

	return c.clientPool[rand.Intn(c.poolSize)], nil
}
