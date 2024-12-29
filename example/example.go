package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"fmt"
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
	conn,err := net.Dial("tcp", "1.1.1.1:80")
	if err != nil {
		slog.Error("Instantiate failed", "err", err)
		return
	}
	f,err := conn.(*net.TCPConn).File()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(f.Fd())
	return
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	_, err = wazero_net.InitFuncExport(r).Instantiate(ctx)
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
	// _, err = r.InstantiateWithConfig(ctx, netWasm, conf)
	// if err != nil {
	// 	log.Panicln(err)
	// }
}
