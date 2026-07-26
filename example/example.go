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
	ctx := context.Background()
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

	// fsPath := os.Args[2]
	// fmt.Println(fsPath)
	// fsConfig := wazero.NewFSConfig()
	// mountFs := `/` + strings.ReplaceAll(strings.ReplaceAll(fsPath, ":", ""), `\`, "/")
	// fsConfig = fsConfig.WithDirMount(fsPath, mountFs)
	// fmt.Println("mount fs", mountFs)

	conf := wazero.NewModuleConfig().WithStartFunctions("_initialize").
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().WithEnv("testenv", "envvalue")
	// WithFSConfig(fsConfig)

	if os.Args[1] == "net" {
		netWasm, err := os.ReadFile("net.wasm")
		if err != nil {
			log.Panicln(err)
		}
		cm, err := r.CompileModule(context.Background(), netWasm)
		if err != nil {
			log.Panicln(err)
		}

		fmt.Printf("cm %+v\n", cm)
		netMod, err := r.InstantiateModule(ctx, cm, conf)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Printf("dial %p\n", netMod)
		netMod.ExportedFunction("dial").Call(context.Background())

		conf2 := wazero.NewModuleConfig().WithStartFunctions("_initialize").
			WithStdout(os.Stdout).
			WithStderr(os.Stderr).
			WithStdin(os.Stdin).
			WithRandSource(rand.Reader).
			WithSysNanosleep().
			WithSysNanotime().
			WithSysWalltime().WithName("main2")

		netMod2, err := r.InstantiateModule(ctx, cm, conf2)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Printf("dial %p\n", netMod2)
		nCtx := context.WithValue(ctx, "test", "value1111")
		netMod2.ExportedFunction("dial").Call(nCtx)
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
		fmt.Printf("InstantiateModule11 %p\n", httpMod)
		malloc := httpMod.ExportedFunction("malloc")
		free := httpMod.ExportedFunction("free")

		https := httpMod.ExportedFunction("https_get")
		s := time.Now()
		url := "http://www.google.com"
		results, err := malloc.Call(ctx, uint64(len(url)))
		if err != nil {
			log.Panicln(err)
		}
		defer free.Call(ctx, results[0])
		httpMod.Memory().Write(uint32(results[0]), []byte(url))

		_, err = https.Call(ctx, results[0], uint64(len(url)))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(time.Since(s))
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
	} else if os.Args[1] == "panic" {
		panicWasm, err := os.ReadFile("panic.wasm")
		if err != nil {
			log.Panicln(err)
		}
		cm, err := r.CompileModule(context.Background(), panicWasm)
		if err != nil {
			log.Panicln(err)
		}
		mod, err := r.InstantiateModule(ctx, cm, conf)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(mod.ExportedFunctionDefinitions())

		_, err = mod.ExportedFunction("panic_test").Call(ctx)
		if err != nil {
			fmt.Println("call panic test")
			fmt.Println(err)
		}
	} else if os.Args[1] == "webdav" {
		panicWasm, err := os.ReadFile("webdav.wasm")
		if err != nil {
			log.Panicln(err)
		}
		cm, err := r.CompileModule(context.Background(), panicWasm)
		if err != nil {
			log.Panicln(err)
		}
		mod, err := r.InstantiateModule(ctx, cm, conf)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(mod.ExportedFunctionDefinitions())

		_, err = mod.ExportedFunction("webdav_connect").Call(ctx)
		if err != nil {
			fmt.Println("call panic test")
			fmt.Println(err)
		}
	} else if os.Args[1] == "localfile" {
		panicWasm, err := os.ReadFile("localfile.wasm")
		if err != nil {
			log.Panicln(err)
		}
		cm, err := r.CompileModule(context.Background(), panicWasm)
		if err != nil {
			log.Panicln(err)
		}
		mod, err := r.InstantiateModule(ctx, cm, conf)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(mod.ExportedFunctionDefinitions())
		malloc := mod.ExportedFunction("malloc")
		free := mod.ExportedFunction("free")

		dir := os.Args[3]

		results, err := malloc.Call(ctx, uint64(len(dir)))
		if err != nil {
			log.Panicln(err)
		}
		defer free.Call(ctx, results[0])
		mod.Memory().Write(uint32(results[0]), []byte(dir))
		_, err = mod.ExportedFunction("localfile").Call(ctx, results[0], uint64(len(dir)))
		if err != nil {
			fmt.Println("call panic test")
			fmt.Println(err)
		}
	}
}
