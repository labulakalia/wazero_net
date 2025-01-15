package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/labulakalia/wazero_net"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed http.wasm
var httpWasm []byte

// //go:embed net.wasm
// var netWasm []byte

func main() {

	ctx := context.Background()
	now := time.Now()
	// add wasm cache
	cacheDir := os.TempDir()
	os.MkdirAll(cacheDir, 0755)
	fmt.Println("cache dir", cacheDir)
	cache, err := wazero.NewCompilationCacheWithDir(cacheDir)
	if err != nil {
		log.Panicln(err)
	}
	rConfig := wazero.NewRuntimeConfig().
		WithCompilationCache(cache)

	r := wazero.NewRuntimeWithConfig(ctx, rConfig)
	defer r.Close(ctx)
	_, err = wazero_net.InitFuncExport(r).Instantiate(ctx)
	if err != nil {
		slog.Error("Instantiate failed", "err", err)
		return
	}

	defer cache.Close(ctx)

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
	cm, err := r.CompileModule(context.Background(), httpWasm)
	if err != nil {
		log.Panicln(err)
	}

	httpsMod, err := r.InstantiateModule(ctx, cm, conf)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("start init module ok ", time.Now().Sub(now))
	malloc := httpsMod.ExportedFunction("malloc")

	url := "https://httpbin.org/get"
	result, err := malloc.Call(ctx, uint64(len(url)))
	if err != nil {
		log.Fatalln("malloc", err)
	}
	httpsMod.Memory().Write(uint32(result[0]), []byte(url))

	_, err = httpsMod.ExportedFunction("https_get").Call(ctx, result[0], uint64(len(url)))
	if err != nil {
		log.Fatalln("https get", err)
	}

	// call 2
	httpsMod, err = r.InstantiateModule(ctx, cm, conf)
	if err != nil {
		log.Panicln(err)
	}

	malloc = httpsMod.ExportedFunction("malloc")

	url = "https://httpbin.org/get"
	result, err = malloc.Call(ctx, uint64(len(url)))
	if err != nil {
		log.Fatalln("malloc", err)
	}
	httpsMod.Memory().Write(uint32(result[0]), []byte(url))

	_, err = httpsMod.ExportedFunction("https_get").Call(ctx, result[0], uint64(len(url)))
	if err != nil {
		log.Fatalln("https get", err)
	}

	// netMod, err := r.InstantiateWithConfig(ctx, netWasm, conf)
	// if err != nil {
	// 	log.Panicln(err)
	// }
	// _, err = netMod.ExportedFunction("net_dial").Call(ctx)
	// if err != nil {
	// 	log.Panicln(err)
	// }
}
