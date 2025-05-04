//go:build wasip1

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/labulakalia/wazero_net/util"
	_ "github.com/labulakalia/wazero_net/wasi/http"
)

//go:wasmexport https_get
func https_get(urlPtr, length uint64) {
	// log(fmt.Sprintf("ptr %d len %d", urlPtr, length))
	geturl := util.PtrToString(uint32(urlPtr), uint32(length))
	resp, err := http.Get(geturl)

	if err != nil {
		log.Panicln("get failed", err)
	}
	fmt.Println("resp status code", resp.StatusCode)
	fmt.Println("resp header", resp.Header)
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln("readall failed", err)
	}
	fmt.Println("resp body\n", string(respData))
}

func main() {}
