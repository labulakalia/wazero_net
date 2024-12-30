package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"log"
	"log/slog"
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

	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	_, err := wazero_net.InitFuncExport(r).Instantiate(ctx)
	if err != nil {
		slog.Error("Instantiate failed", "err", err)
		return
	}
	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	conf := wazero.NewModuleConfig().
		WithStartFunctions("_initialize").
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime()
	httpsMod, err := r.InstantiateWithConfig(ctx, httpWasm, conf)
	if err != nil {
		log.Panicln(err)
	}
	_, err = httpsMod.ExportedFunction("https_get").Call(ctx)
	if err != nil {
		log.Panicln(err)
	}

	netMod, err := r.InstantiateWithConfig(ctx, netWasm, conf)
	if err != nil {
		log.Panicln(err)
	}
	_, err = netMod.ExportedFunction("net_dial").Call(ctx)
	if err != nil {
		log.Panicln(err)
	}
}
