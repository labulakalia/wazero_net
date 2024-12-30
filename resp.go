package wazero_net

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Response struct {
	StatusCode int    // e.g. 200
	Proto      string // e.g. "HTTP/1.0"
	ProtoMajor int    // e.g. 1
	ProtoMinor int    // e.g. 0

	Header http.Header

	Body []byte

	ContentLength int64
}

func (r *Response) ToHttpResponse() *http.Response {
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

type Request struct {
	Method string

	URL *url.URL

	Proto      string // "HTTP/1.0"
	ProtoMajor int    // 1
	ProtoMinor int    // 0

	Header http.Header

	Body []byte

	ContentLength int64

	TransferEncoding []string

	Close bool

	Host string

	Form url.Values

	PostForm url.Values

	MultipartForm *multipart.Form

	RemoteAddr string

	RequestURI string

	Pattern string
}

func (r *Request) ToHttpRequest() (req *http.Request) {
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

func (r *Request) ParseHttpRequest(req *http.Request) error {
	r.Method = req.Method
	r.URL = req.URL
	r.Proto = req.Proto
	r.ProtoMajor = req.ProtoMajor
	r.ProtoMinor = req.ProtoMinor
	r.Header = req.Header
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
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
	return nil
}
