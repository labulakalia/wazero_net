//go:build wasip1

package http

import (
	"net/http"
	"time"
)

func init() {
	http.DefaultTransport = &Transport{}
	http.DefaultClient.Transport = &Transport{}
}

type Transport struct {
	timeout time.Duration
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	defaultTimeout := http.DefaultClient.Timeout
	return Do(req, defaultTimeout)
}
