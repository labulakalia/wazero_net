package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"runtime"

	_ "github.com/labulakalia/wazero_net/wasi/malloc"
	"github.com/labulakalia/wazero_net/wasi/net"
)

//go:wasmexport net_dial
func net_dial() {

	fmt.Println(runtime.GOMAXPROCS(10))
	// fmt.Println()
	slog.SetLogLoggerLevel(slog.LevelDebug)
	lis, err := net.Listen("tcp", "0.0.0.0:19971")
	if err != nil {
		slog.Error("listen failed", "err", err)
		return
	}
	go startListen(lis)

	conn, err := net.Dial("tcp", "127.0.0.1:19971")
	if err != nil {
		slog.Error("dial failed", "err", err)
		return
	}
	data := make([]byte, 1024)
	for i := 0; i < 10; i++ {
		slog.Info("start write")
		n, err := conn.Write([]byte("hello"))
		if err != nil {
			slog.Error("write failed", "err", err)
		}
		slog.Info("write success", "n", n)
		// time.Sleep(time.Millisecond) // must exist
		// runtime.Gosched()
		n, err = conn.Read(data)
		if err != nil {
			slog.Error("write failed", "err", err)
		}
		slog.Info("read", "n", n)

	}
}

func startListen(lis *net.Listener) {

	for {
		conn, err := lis.Accept()
		if err != nil {
			slog.Error("accept failed", "err", err)
			continue
		}
		data := make([]byte, 1024)
		for {
			slog.Info("start read")
			n, err := conn.Read(data)
			if err != nil {
				if errors.Is(err, io.EOF) {
					slog.Info("read success")
					return
				}
				slog.Error("read failed", "err", err)
				return
			}
			slog.Info("read success", "data", string(data[:n]))

			wn, err := conn.Write(data[:n])
			if err != nil {
				slog.Error("write failed", "err", err)
				return
			}
			if wn != n {
				slog.Error("read count not equal", "rn", wn, "n", n)
				break
			}
			// runtime.Gosched()
			// time.Sleep(time.Millisecond) // must exist
		}
	}
}

//go:wasmexport dial
func dial() {
	conn, err := net.Dial("tcp", "127.0.0.1:2121")
	if err != nil {
		slog.Error("wasm dial failed", "err", err)
		return
	}
	fmt.Println(conn.RemoteAddr())
	conn.Close()

	conn, err = net.Dial("tcp", "127.0.0.1:9999")
	if err == nil {
		log.Panicln("err should not nil")
	}
	fmt.Println("dial err", err)

}

func main() {}
