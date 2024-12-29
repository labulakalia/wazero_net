package main

import (
	"log/slog"
	"wazero_net/wasi/net"
)


func main(){
	conn, err := net.Dial("tcp", "1.1.1.1:80")
	if err != nil {
		slog.Error("dial failed", "err", err)
		return
	}
	n,err := conn.Write([]byte("hello"))
	if err != nil {
		slog.Error("write failed", "err",err)
	}
	slog.Info("write", "n",n)
}
