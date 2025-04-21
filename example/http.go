//go:build wasip1

package main

import (
	"net/http"
	"runtime"
	"unsafe"

	"github.com/labulakalia/wazero_net/util"

	wasihttp "github.com/labulakalia/wazero_net/wasi/http"
)

// _ "github.com/labulakalia/wazero_net/wasi/http"
//
//go:wasmimport env log
func _log(ptr, size uint32)

func stringToPtr(s string) (uint32, uint32) {
	ptr := unsafe.Pointer(unsafe.StringData(s))
	return uint32(uintptr(ptr)), uint32(len(s))
}

func log(message string) {
	ptr, size := stringToPtr(message)
	_log(ptr, size)
	runtime.KeepAlive(message) // keep message alive until ptr is no longer needed.
}

//go:wasmexport https_get
func https_get(urlPtr, length uint64) {
	// log(fmt.Sprintf("ptr %d len %d", urlPtr, length))
	geturl := util.PtrToString(uint32(urlPtr), uint32(length))
	log(geturl)
	http.DefaultTransport = &wasihttp.Transport{}
	_, err := http.Get("https://httpbin.org/get")
	log(err.Error())

	// slog.Info("http get", "ptr", urlPtr, "length", length)
	// slog.Info("get url", "url", geturl)
	// u, err := url.Parse(geturl)
	// if err != nil {
	// 	log.Panicln("parse url failed", err)
	// }
	// u = u
	// resp, err := http.Get(u.String())

	// if err != nil {
	// 	log.Panicln("do failed", err)
	// }
	// fmt.Println("resp status code", resp.StatusCode)
	// fmt.Println("resp header", resp.Header)
	// fmt.Println("resp body", resp.Body)
}

func main() {}
