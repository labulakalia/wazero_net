package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cloudsoda/go-smb2"

	wasi_net "github.com/labulakalia/wazero_net/wasi/net"
)

func main() {}

//go:wasmexport samba_connect
func samba_connect() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	smbDialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     "user",
			Password: "passwd",
		},
	}
	slog.Info("tcp dial")
	conn, err := wasi_net.Dial("tcp", "127.0.0.1:445")
	if err != nil {
		slog.Error("failed to dial", "error", err)
		return
	}
	slog.Info("smb dial")

	smbSession, err := smbDialer.DialConn(context.Background(), conn, "127.0.0.1:445")
	if err != nil {
		slog.Error("failed to dial", "error", err)
		return
	}
	smbSession = smbSession
	fmt.Println(smbSession.ListSharenames())
}
