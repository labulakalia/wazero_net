package wazero_net

import (
	"bytes"
	"io"
	"net/http"

	"github.com/labulakalia/wazero_net/model"
)

func toHttpResponse(r *model.Response) *http.Response {
	return &http.Response{
		Status:        http.StatusText(r.StatusCode),
		StatusCode:    r.StatusCode,
		Proto:         r.Proto,
		ProtoMajor:    r.ProtoMajor,
		ProtoMinor:    r.ProtoMinor,
		Header:        r.Header,
		Body:          io.NopCloser(bytes.NewBuffer(r.Body)),
		ContentLength: r.ContentLength,
	}
}

func toHttpRequest(r *model.Request) (req *http.Request) {
	req = &http.Request{}
	req.Method = r.Method
	req.URL = r.URL
	req.Proto = r.Proto
	req.ProtoMajor = r.ProtoMajor
	req.ProtoMinor = r.ProtoMinor
	req.Header = r.Header
	if r.Body != nil {
		req.Body = io.NopCloser(bytes.NewBuffer(r.Body))
	}
	req.ContentLength = r.ContentLength
	req.TransferEncoding = r.TransferEncoding
	req.Host = r.Host
	req.Form = r.Form
	req.MultipartForm = r.MultipartForm
	req.PostForm = r.PostForm
	req.RemoteAddr = r.RemoteAddr
	req.RequestURI = r.RequestURI
	return req
}

func toModelRequest(req *http.Request) (*model.Request, error) {
	r := &model.Request{}
	r.Method = req.Method
	r.URL = req.URL
	r.Proto = req.Proto
	r.ProtoMajor = req.ProtoMajor
	r.ProtoMinor = req.ProtoMinor
	r.Header = req.Header
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		r.Body = bodyBytes
	}
	r.ContentLength = req.ContentLength
	r.TransferEncoding = req.TransferEncoding
	r.Host = req.Host
	r.Form = req.Form
	r.MultipartForm = req.MultipartForm
	r.PostForm = req.PostForm
	r.RemoteAddr = req.RemoteAddr
	r.RequestURI = req.RequestURI
	return r, nil
}
