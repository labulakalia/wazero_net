//go:build wasip1

package main

import (
	"fmt"
	"log"
	"log/slog"

	"net/http"
	"net/url"
	"os"
	"syscall"

	_ "github.com/labulakalia/plugin_api"
	"github.com/labulakalia/wazero_net/model"
	"github.com/labulakalia/wazero_net/util"
	wasihttp "github.com/labulakalia/wazero_net/wasi/http"
)

var isRun = false

//go:wasmexport https_get
func https_get(urlPtr, length uint64) {
	entries, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			panic(err)
		}
		if info.Sys() != nil {
			fmt.Println(info.Sys().(*syscall.Stat_t).Atime / 1000000)
		}
	}
	return
	client := http.Client{}
	client = client
	fmt.Println("is run", isRun)
	isRun = true

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
