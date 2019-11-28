package http

import (
	"net"
	"net/http"
	"time"
)

// NewClient returns new http.Client optimised for high amount of outgoing requests,
// using increased amount of maximum idle connections per host in client's pool.
// connCount controls the amount of connections in the pool.
func NewClient(connCount uint32) *http.Client {
	// Taken from http://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		DisableKeepAlives:     false,
		MaxIdleConns:          int(connCount),
		MaxIdleConnsPerHost:   int(connCount),
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport: transport,
	}
}
