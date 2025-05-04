package main

import (
	"errors"
	"io"
	"log"
	"log/slog"
	"net"

	// "net/http"

	_ "github.com/labulakalia/wazero_net/wasi/http"
	wasinet "github.com/labulakalia/wazero_net/wasi/net"
)

//go:wasmexport net_dial
func net_dial() {
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
	conn.Close()
}

func startListen(lis net.Listener) {
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

		}
	}
}

var conn net.Conn

//go:wasmexport dial
func dial() {
	conn, err := wasinet.Dial("tcp", "1.1.1.1:80")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(conn.RemoteAddr())
	log.Println(conn.LocalAddr())
}

func main() {}
