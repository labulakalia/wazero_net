package http

import "net/http"

func init() {
	http.DefaultTransport = &Transport{}
	http.DefaultClient.Transport = &Transport{}
}

type Transport struct {
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return Do(req)
}
