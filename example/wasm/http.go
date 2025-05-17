//go:build wasip1

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/labulakalia/wazero_net/util"
	_ "github.com/labulakalia/wazero_net/wasi/http"
)

//go:wasmexport https_get1
func https_get1(urlPtr, length uint64) {
	// log(fmt.Sprintf("ptr %d len %d", urlPtr, length))
	geturl := util.PtrToString(uint32(urlPtr), uint32(length))
	resp, err := http.Get(geturl)

	if err != nil {
		log.Panicln("get failed", err)
	}
	fmt.Println("resp status code", resp.StatusCode)
	fmt.Println("resp header", resp.Header)
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln("readall failed", err)
	}
	fmt.Println("resp body\n", string(respData))
}

//go:wasmexport https_get2
func https_get2(urlPtr, length uint64) {
	// log(fmt.Sprintf("ptr %d len %d", urlPtr, length))
	geturl := util.PtrToString(uint32(urlPtr), uint32(length))
	resp, err := http.Get(geturl)

	if err != nil {
		log.Panicln("get failed", err)
	}
	fmt.Println("resp status code", resp.StatusCode)
	fmt.Println("resp header", resp.Header)
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln("readall failed", err)
	}
	fmt.Println("resp body\n", string(respData))
}

//go:wasmexport http_redirect
func http_redirect() {
	client :=
		http.DefaultClient
	r, _ := http.NewRequest(http.MethodGet, "http://192.168.123.29:5244/dav/tianyi/%E6%88%91%E7%9A%84%E8%A7%86%E9%A2%91/%E7%BB%9D%E5%AF%B9%E6%9D%83%E5%8A%9B%5B%E7%AE%80%E7%B9%81%E8%8B%B1%E5%AD%97%E5%B9%95%5D.Absolute.Power.1997.EUR.1080p.BluRay.x265.10bit.DTS-SONYHD/Absolute.Power.1997.EUR.1080p.BluRay.x265.10bit.DTS-SONYHD.mkv", nil)
	r.Header.Set("Authorization", "Basic YWRtaW46MTA5MDk3")
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Location())
}
func main() {}
