package main

import (
	"context"
	"io"
	"log/slog"
	"net"

	"net/http"
	// _ "wazero_net/wasi/http"
	wazero_net "wazero_net/wasi/net"
)

func main() {
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {

		conn,err := wazero_net.Dial(network, addr)
		if err != nil {
			return nil,err
		}
		slog.Info("start dial", "addr",conn.RemoteAddr())
		return conn,nil
	}
	req,err := http.NewRequest(http.MethodGet, "http://192.168.123.53:8000", nil)
	if err != nil {
		slog.Error("new request failed", "err", err)
		return
	}
	resp,err := http.DefaultClient.Do(req)
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
