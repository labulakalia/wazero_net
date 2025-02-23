package http

//go:wasmimport net client_do
//go:noescape
func client_do(reqPtr, reqLen, respPtr, respLen uint64) uint64
