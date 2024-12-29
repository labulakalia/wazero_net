package http

//go:wasmimport net round_trip
//go:noescape
func round_trip(reqPtr, reqLen,respLenPtr uint64) uint64


//go:wasmimport net read_resp
//go:noescape
func read_resp(dataPtr, dataLen uint64) uint64
