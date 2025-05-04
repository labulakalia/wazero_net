// TINYGO: The following is copied and modified from Go 1.21.4 official implementation.

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// HTTP client implementation. See RFC 7230 through 7235.
//
// This is the low-level Transport implementation of RoundTripper.
// The high-level interface is in client.go.

package http

import (
	"io"
	"net/http"
	"unsafe"
)

type readTrackingBody struct {
	io.ReadCloser
	didRead  bool
	didClose bool
}

type Transport struct{}

var DefaultTransport http.RoundTripper = &Transport{}

// roundTrip implements a RoundTripper over HTTP.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// return roundTrip(req)
	_req := (*Request)(unsafe.Pointer(req))
	_resp, err := roundTrip(_req)
	if err != nil {
		return nil, err
	}
	resp := (*http.Response)(unsafe.Pointer(_resp))
	return resp, nil
}

func init() {
	http.DefaultClient.Transport = &Transport{}
}
