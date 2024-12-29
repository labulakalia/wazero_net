package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"errors"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"wazero_net"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed http.wasm
var httpWasm []byte

//go:embed net.wasm
var netWasm []byte

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	// go startListen()
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	_, err := wazero_net.InitFuncExport(r).Instantiate(ctx)
	if err != nil {
		slog.Error("Instantiate failed", "err", err)
		return
	}
	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	conf := wazero.NewModuleConfig().

		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime()
	_, err = r.InstantiateWithConfig(ctx, httpWasm, conf)
	if err != nil {
		log.Panicln(err)
	}
	_, err = r.InstantiateWithConfig(ctx, netWasm, conf)
	if err != nil {
		log.Panicln(err)
	}
}

func startListen() {
	lis, err := net.Listen("tcp", "0.0.0.0:19971")
	if err != nil {
		slog.Error("listen failed", "err", err)
		return
	}
	slog.Info("start listen", "addr", lis.Addr())
	for {
		conn, err := lis.Accept()
		if err != nil {
			slog.Error("accept failed", "err", err)
			continue
		}
		data := make([]byte, 1024)
		for {
			slog.Info("wait read")
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
