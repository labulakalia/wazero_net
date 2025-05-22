//go:build wasip1

package main

import (
	"log/slog"
	"net/http"
	"time"

	_ "github.com/labulakalia/wazero_net/wasi/http"
	"github.com/medianexapp/gowebdav"
	// "github.com/medianexapp/plugin_api/httpclient"
)

//go:wasmexport webdav_connect
func webdav_connect() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	client := gowebdav.NewClient("", "", "")
	client.SetClientDo(http.DefaultClient.Do)
	err := client.Connect()
	if err != nil {
		slog.Error("connect failed", "err", err)
		return
	}
	client.SetTimeout(time.Second * 30)
	dirs, err := client.ReadDir("/")
	if err != nil {
		slog.Error("read dir failed", "err", err)
		return
	}
	for _, dir := range dirs {
		slog.Info("read dir", "dir", dir.Name())
		if dir.IsDir() {

			dirs, err := client.ReadDir("/" + dir.Name())
			if err != nil {
				slog.Error("read dir failed", "err", err)
				return
			}
			for _, dir := range dirs {
				slog.Info("read dir", "dir", dir.Name())
			}
		}
	}
}

func main() {}
