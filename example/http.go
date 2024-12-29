package main

import (
	"io"
	"log/slog"
	"net/http"
	_ "wazero_net/wasi/http"
)

func main() {
	resp,err := http.Get("https://www.baidu.com")
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
