package main

import (
	"fmt"
	"log"
	"log/slog"

	"net/url"

	"github.com/labulakalia/wazero_net/model"
	"github.com/labulakalia/wazero_net/util"
	wasihttp "github.com/labulakalia/wazero_net/wasi/http"
)

//go:wasmexport https_get
func https_get(urlPtr, length uint64) {
	geturl := util.PtrToString(uint32(urlPtr), uint32(length))
	fmt.Println("start http get")
	slog.Info("get url", "url", geturl)
	u, err := url.Parse(geturl)
	if err != nil {
		log.Panicln("parse url failed", err)
	}
	slog.Info("get url", "url", u)
	u = u
	resp, err := wasihttp.Do(&model.Request{
		Method: "GET",
		URL:    u,
	})
	slog.Info("get url", "url", err)
	if err != nil {
		log.Panicln("do failed", err)
	}
	fmt.Println("resp", resp)
}

func main() {}
