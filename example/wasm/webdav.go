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
	client := gowebdav.NewClient("http://192.168.123.29:5244/dav", "admin", "109097")
	client.SetClientDo(http.DefaultClient.Do)
	err := client.Connect()
	if err != nil {
		slog.Error("connect failed", "err", err)
		// return err
	}
	req, err := client.GetPathRequest("/dav/tianyi/%E6%88%91%E7%9A%84%E8%A7%86%E9%A2%91/%E7%BB%9D%E5%AF%B9%E6%9D%83%E5%8A%9B%5B%E7%AE%80%E7%B9%81%E8%8B%B1%E5%AD%97%E5%B9%95%5D.Absolute.Power.1997.EUR.1080p.BluRay.x265.10bit.DTS-SONYHD/Absolute.Power.1997.EUR.1080p.BluRay.x265.10bit.DTS-SONYHD.mkv")

	fmt.Println(err)
	fmt.Println(req.Header)
	fmt.Println(req.URL)
}

func main() {}
