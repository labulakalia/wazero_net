package main

import (
	"io"
	"log/slog"

	"net/http"

	_ "github.com/labulakalia/wazero_net/wasi/http"
)

//go:wasmexport https_get
func https_get() {
	req, err := http.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
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
	slog.Info("get resp", "data", string(bytes))
}

func main() {}
