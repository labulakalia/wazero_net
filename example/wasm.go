package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
	_ "wazero_net/wasi/http"
	"wazero_net/wasi/net"
)

// NOTE: multi goroutine can not scheduler on wasm,
func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	for {
		httpsGet()
		time.Sleep(time.Millisecond*10)
	}
	// dialTls()
	// netDial()
}

func dialTls() {

	conn, err := net.DialTls("tcp4", "www.baidu.com:443")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(conn.RemoteAddr())
	conn.Close()

}

func netDial() {
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

//go:wasmexport httpsGet
func httpsGet() {
	os.Stdout.WriteString("start get")
	resp, err := http.Get("http://192.168.123.53:8000")
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
