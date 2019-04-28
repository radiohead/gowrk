package types

import (
	"net/http"
)

// HTTPClient implements the request sending method.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
