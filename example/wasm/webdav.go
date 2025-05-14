//go:build wasip1

package main

import (
	"fmt"
	"log/slog"
	"net/http"

	_ "github.com/labulakalia/wazero_net/wasi/http"
	"github.com/medianexapp/gowebdav"
	// "github.com/medianexapp/plugin_api/httpclient"
)

//go:wasmexport webdav_connect
func webdav_connect() {
	// c := httpclient.NewClient()
	// log(fmt.Sprintf("ptr %d len %d", urlPtr, length))
	client := gowebdav.NewClient("https://127.0.0.1:8080", "user", "passwd")
	client.SetClientDo(http.DefaultClient.Do)
	err := client.Connect()
	if err != nil {
		slog.Error("connect failed", "err", err)
		// return err
	}
	fmt.Println(client.GetPathRequest("/"))
}

func main() {}
