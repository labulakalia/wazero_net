package main

import (
	"fmt"
	"log/slog"
	"wazero_net/wasi/net"
	_ "wazero_net/wasi/net"
)

//go:wasmexport httpGet
func httpGet() {

}

// NOTE: multi goroutine can not scheduler on wasm,
func main() {
	slog.Info("dial", "addr", "127.0.0.1:19971")
	conn, err := net.Dial("tcp", "127.0.0.1:19971")
	if err != nil {
		slog.Error("dial failed", "err", err)
		return
	}
	data := make([]byte, 1024)
	for i := 0; i < 100; i++ {
		slog.Info("start write")
		str := []byte(fmt.Sprintf("data data %d", i))
		n, err := conn.Write(str)
		if err != nil {
			slog.Error("write failed", "err", err)
			break
		}
		slog.Info("write success", "count", n)
		rn, err := conn.Read(data)
		if err != nil {
			slog.Error("write failed", "err", err)
			break
		}
		if rn != n {
			slog.Error("read count not equal", "rn", rn, "n", n)
			break
		}

	}
	conn.Close()
}
