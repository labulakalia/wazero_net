package http

//go:wasmimport net client_do
//go:noescape
func client_do(reqPtr, reqLen uint64) uint64
