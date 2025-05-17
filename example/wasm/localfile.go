//go:build wasip1

package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/labulakalia/wazero_net/util"
	_ "github.com/labulakalia/wazero_net/wasi/http"
	// "github.com/medianexapp/plugin_api/httpclient"
)

//go:wasmexport localfile
func localfile(urlPtr, length uint64) {
	dir := util.PtrToString(uint32(urlPtr), uint32(length))
	dirs, err := os.ReadDir(dir)
	if err != nil {
		slog.Error("read failed", "err", err)
		return
	}
	fmt.Println(os.Getwd())
	for _, dir := range dirs {
		fmt.Println("dir %s", dir.Name())
	}

}

func main() {}
