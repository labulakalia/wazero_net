//go:build wasip1

package main

import (
	"fmt"
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
		log.Panicln("do failed", err)
	}
	fmt.Println("resp status code", resp.StatusCode)
	fmt.Println("resp header", resp.Header)
	fmt.Println("resp body", resp.Body)
}

func main() {}
