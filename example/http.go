package main

import (
	"fmt"
	"io"
	"log/slog"

	"net/http"

	"github.com/labulakalia/wazero_net/util"
	_ "github.com/labulakalia/wazero_net/wasi/http"
)

//go:wasmexport https_get
func https_get(urlPtr, length uint64) {
	url := util.PtrToString(uint32(urlPtr), uint32(length))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("new request failed", "err", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("http get failed", "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("http status code failed", "status", resp.Status)
		return
	}
	slog.Info("http resp", "header", resp.Header)
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("read failed", "err", err)
		return
	}
	fmt.Println(string(bytes))
}

func main() {}
