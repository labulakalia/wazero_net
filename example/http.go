//go:build wasip1

package main

import (
	"fmt"
	"log"
	"log/slog"

	"net/http"
	"net/url"

	"github.com/labulakalia/wazero_net/util"
	wasihttp "github.com/labulakalia/wazero_net/wasi/http"
)

//go:wasmexport https_get
func https_get(urlPtr, length uint64) {
	geturl := util.PtrToString(uint32(urlPtr), uint32(length))

	slog.Info("get url", "url", geturl)
	u, err := url.Parse(geturl)
	if err != nil {
		log.Panicln("parse url failed", err)
	}

	resp, err := wasihttp.Do(&http.Request{
		Method: "GET",
		URL:    u,
	})

	if err != nil {
		log.Panicln("do failed", err)
	}
	fmt.Println("status code", resp.StatusCode)
	fmt.Println("body", resp.Body)
}

func main() {}
