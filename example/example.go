package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"log"
	"log/slog"
	"os"

	"github.com/labulakalia/wazero_net"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {

	ctx := context.Background()

	// add wasm cache
	cacheDir := os.TempDir()
	os.MkdirAll(cacheDir, 0755)

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

	if os.Args[1] == "net" {
		netWasm, err := os.ReadFile("net.wasm")
		if err != nil {
			log.Panicln(err)
		}
		cm, err := r.CompileModule(context.Background(), netWasm)
		if err != nil {
			log.Panicln(err)
		}
		netMod, err := r.InstantiateModule(ctx, cm, conf)
		if err != nil {
			log.Panicln(err)
		}
		netMod.ExportedFunction("net_dial").Call(context.Background())
		// netMod.ExportedFunction("net_dial").Call(context.Background())
	} else if os.Args[1] == "http" {
		netWasm, err := os.ReadFile("http.wasm")
		if err != nil {
			log.Panicln(err)
		}
		cm, err := r.CompileModule(context.Background(), netWasm)
		if err != nil {
			log.Panicln(err)
		}

		httpMod, err := r.InstantiateModule(ctx, cm, conf)
		if err != nil {
			log.Panicln(err)
		}
		malloc := httpMod.ExportedFunction("malloc")
		url := "https://httpbin.org/get"
		result, err := malloc.Call(ctx, uint64(len(url)))
		if err != nil {
			log.Fatalln("malloc", err)
		}
		httpMod.Memory().Write(uint32(result[0]), []byte(url))

		_, err = httpMod.ExportedFunction("https_get").Call(ctx, result[0], uint64(len(url)))
		if err != nil {
			log.Fatalln("https get", err)
		}
	} else if os.Args[1] == "ftp" {
		ftpWasm, err := os.ReadFile("ftp.wasm")
		if err != nil {
			log.Panicln(err)
		}
		cm, err := r.CompileModule(context.Background(), ftpWasm)
		if err != nil {
			log.Panicln(err)
		}
		netMod, err := r.InstantiateModule(ctx, cm, conf)
		if err != nil {
			log.Panicln(err)
		}

		netMod.ExportedFunction("ftp_connect").Call(context.Background())
	} else if os.Args[1] == "sftp" {
		sftpWasm, err := os.ReadFile("sftp.wasm")
		if err != nil {
			log.Panicln(err)
		}
		cm, err := r.CompileModule(context.Background(), sftpWasm)
		if err != nil {
			log.Panicln(err)
		}
		netMod, err := r.InstantiateModule(ctx, cm, conf)
		if err != nil {
			log.Panicln(err)
		}

		netMod.ExportedFunction("sftp_connect").Call(context.Background())
	} else if os.Args[1] == "smb" {
		smbWasm, err := os.ReadFile("smb.wasm")
		if err != nil {
			log.Panicln(err)
		}
		cm, err := r.CompileModule(context.Background(), smbWasm)
		if err != nil {
			log.Panicln(err)
		}
		mod, err := r.InstantiateModule(ctx, cm, conf)
		if err != nil {
			log.Panicln(err)
		}

		mod.ExportedFunction("smb_connect").Call(context.Background())
	}
}
