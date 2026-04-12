//go:build wasip1

package main

import (
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/labulakalia/wazero_net/util"
	_ "github.com/labulakalia/wazero_net/wasi/http"
	"github.com/medianexapp/plugin_api/httpclient"
)

//go:wasmexport https_get
func https_get(urlPtr, length uint64) {
	geturl := util.PtrToString(uint32(urlPtr), uint32(length))
	slog.Info("start http get", "url", geturl)
	client := httpclient.NewClient(httpclient.WithTimeout(time.Second * 3))
	resp, err := client.Get(geturl)
	if err != nil {
		slog.Error("client get failed", "err", err)
		return
	}
	fmt.Println("resp status code", resp.StatusCode)
	fmt.Println("resp header", resp.Header)
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("read failed", "err", err)
		return
	}
	fmt.Println("resp body\n", string(respData))
}

func main() {}
