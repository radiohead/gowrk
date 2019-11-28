package http

import (
	"errors"
	"net/http"
	"sync/atomic"
)

// ErrEmptyPool is returned if the client pool is zero-size.
var ErrEmptyPool = errors.New("pool of size 0 is not allowed")

//go:generate mockgen -destination=../../test/mocks/http_mocks.go -package=mocks -source=pooled_client.go Client

// Client does a barrlen
type Client interface {
	Do(*http.Request) (*http.Response, error)
}

// PooledClient contains a pool of http.Client instances.
// The client implements a reduced interface of http.Client.
type PooledClient struct {
	poolSize   uint32
	poolIdx    uint32
	clientPool []Client
}

// NewPooledClient returns a new instance of PooledClient, ready to be used.
// poolSize controls the amount of clients in the pool
// and maxIdleConns controls the size of the connection pool of individual clients.
func NewPooledClient(poolSize uint32, maxIdleConns uint32) (*PooledClient, error) {
	clients := make([]Client, poolSize)
	for i := 0; i < int(poolSize); i++ {
		clients[i] = NewClient(maxIdleConns)
	}

	return NewPooledClientWithPool(clients)
}

// NewPooledClientWithPool returns a new instance of PooledClient,
// with client pool set from the pool, which is copied and can be later re-used safely.
func NewPooledClientWithPool(pool []Client) (*PooledClient, error) {
	poolSize := len(pool)
	if poolSize < 1 {
		return nil, ErrEmptyPool
	}

	clientPool := make([]Client, poolSize)
	copy(clientPool, pool)

	return &PooledClient{
		clientPool: clientPool,
		poolSize:   uint32(poolSize),
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

func (c *PooledClient) getClient() (Client, error) {
	if c.poolSize < 1 {
		return nil, ErrEmptyPool
	}

	if c.poolSize == 1 {
		return c.clientPool[0], nil
	}

	// Take the next client in the pool.
	// uint32 overflow resets to 0.
	idx := atomic.AddUint32(&c.poolIdx, 1) % c.poolSize

	return c.clientPool[idx], nil
}
