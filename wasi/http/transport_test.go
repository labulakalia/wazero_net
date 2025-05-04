package http

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"
	"unsafe"
)

func TestTransport(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 9, 10}
	req := &http.Request{
		Method: http.MethodGet,
		Body:   io.NopCloser(bytes.NewBuffer(data)),
	}
	_req := (*Request)(unsafe.Pointer(req))

	t.Log(_req.Method)
	res, err := io.ReadAll(_req.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reflect.DeepEqual(res, data))

}
