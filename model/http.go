package model

import (
	"mime/multipart"

	"net/url"
)

type Response struct {
	StatusCode int    // e.g. 200
	Proto      string // e.g. "HTTP/1.0"
	ProtoMajor int    // e.g. 1
	ProtoMinor int    // e.g. 0

	Header map[string][]string

	Body []byte

	ContentLength int64
}

type Request struct {
	Method string

	URL *url.URL

	Proto      string // "HTTP/1.0"
	ProtoMajor int    // 1
	ProtoMinor int    // 0

	Header map[string][]string

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
