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
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// // add wasm cache
	cacheDir := os.TempDir()
	os.MkdirAll(cacheDir, 0755)

	cache, err := wazero.NewCompilationCacheWithDir(cacheDir)
	if err != nil {
		log.Panicln(err)
	}
	features := api.CoreFeaturesV2.SetEnabled(api.CoreFeatureMutableGlobal, false)
	rConfig := wazero.NewRuntimeConfigCompiler().
		WithCompilationCache(cache).WithDebugInfoEnabled(true).WithCloseOnContextDone(true).WithCoreFeatures(features)

	r := wazero.NewRuntimeWithConfig(ctx, rConfig)
	// defer r.Close(ctx)
	// r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) //
	_, err = wazero_net.InitFuncExport(r).Instantiate(ctx)
	if err != nil {
		slog.Error("Instantiate failed", "err", err)
		return
	}

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	conf := wazero.NewModuleConfig().WithStartFunctions("_initialize").
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
		netMod.ExportedFunction("dial").Call(context.Background())
		// netMod.ExportedFunction("dial").Call(context.Background())
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
		https := httpMod.ExportedFunction("https_get")

		malloc := httpMod.ExportedFunction("malloc")
		free := httpMod.ExportedFunction("free")
		url := "https://httpbin.org/get"
		results, err := malloc.Call(ctx, uint64(len(url)))
		if err != nil {
			log.Panicln(err)
		}
		defer free.Call(ctx, results[0])
		httpMod.Memory().Write(uint32(results[0]), []byte(url))

		_, err = https.Call(ctx, results[0], uint64(len(url)))
		if err != nil {
			log.Fatalln("https get", err)
		}
		httpMod.Close(context.Background())
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

		f := netMod.ExportedFunction("sftp_connect")
		fmt.Println(f.Call(context.Background()))
		// time.Sleep(time.Second * 10)
		fmt.Println(f.Call(context.Background()))

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
		fmt.Println(mod.ExportedFunctionDefinitions())

		_, err = mod.ExportedFunction("smb_connect").Call(ctx)
		if err != nil {
			fmt.Println(err)
		}

	}
}
