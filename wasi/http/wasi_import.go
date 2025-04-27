//go:build wasip1

package http

//go:wasmimport net client_do
//go:noescape
func _client_do(reqPtr, reqLen, respPtr, respLen uint64) uint64
