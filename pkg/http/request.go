package http

import (
	"net/http"
	"net/url"
)

// NewGETRequest returns new http.Request which uses HTTP GET,
// provided URL and an empty body.
func NewGETRequest(url url.URL) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, url.String(), nil)
}
