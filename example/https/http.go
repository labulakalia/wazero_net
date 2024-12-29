package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"

	"net/http"
)

func main() {
	fmt.Println(http.DefaultTransport.RoundTrip)
	return
	// http.DefaultClient.Transport = &http.Transport{
	// 	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
	// 		return net.Dial(network, addr)
	// 	},
	// }
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		slog.Error("http get failed", "err", err)
		return
	}
	ddd,_ := json.Marshal(resp)
	fmt.Println("xx",string(ddd),resp.Status)
	// defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("http status code failed", "status", resp.Status)
		return
	}
	slog.Info("http resp", "header", resp.Header)
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("read failed", "err", err)
		return
	}
	slog.Info("get resp", "data", string(bytes))
	return
	conf := &tls.Config{}
	// //TLS connection
	tlsCon, err := tls.Dial("tcp4", "www.baidu.com:443", conf)
	if err != nil {
		fmt.Println("SSL Error : " + err.Error())
		return
	}

	defer tlsCon.Close()
	fmt.Println(tlsCon.NetConn().RemoteAddr())
	state := tlsCon.ConnectionState()
	fmt.Println("SSL ServerName : " + state.ServerName)
	fmt.Println("SSL Handshake : ", state.HandshakeComplete)
	bb, _ := json.Marshal(state)
	fmt.Println(len(bb))
	// fmt.Println(tlsCon.Handshake())

	request := "GET / HTTP/1.1\r\nHost: www.baidu.com\r\n\r\n"
	n, err := io.WriteString(tlsCon, request)
	if err != nil {
		fmt.Println("SSL Write error :", err.Error(), n)
	}
	data := make([]byte, 65535)
	for {
		fmt.Println("start read")
		n, err = tlsCon.Read(data)
		fmt.Println("err", err)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		fmt.Println(data[n-1], []byte("\n"))
		if data[n-1] == 10 {
			break
		}
		// fmt.Printf("%s", data[:n])
	}
	conn, err := tls.Dial("tcp", "www.baidu.com:443", nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(len([]byte(conn.RemoteAddr().String())))
	// fmt.Println(conn.Write([]byte("hello")))
	conn.Close()
	return
	transport := http.DefaultTransport.(*http.Transport)

	transport.DialTLSContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		slog.Info("dial tls", "network", network, "addr", addr)

		return tls.Dial(network, addr, &tls.Config{})
	}
	resp, err = http.Get("https://www.baidu.com")
	if err != nil {
		slog.Error("http get failed", "err", err)
		return
	}
	fmt.Println(resp.Status)

}
