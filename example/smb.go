package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cloudsoda/go-smb2"

	wasi_net "github.com/labulakalia/wazero_net/wasi/net"
)

func main() {}

//go:wasmexport smb_connect
func smb_connect() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	smbDialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     "labulakalia",
			Password: "109097",
		},
	}
	smbDialer = smbDialer
	addr := "127.0.0.1:445"
	conn, err := wasi_net.Dial("tcp", addr)
	if err != nil {
		slog.Error("failed to dial", "error", err)
		return
	}

	smbSession, err := smbDialer.DialConn(context.Background(), conn, addr)
	if err != nil {
		slog.Error("failed to dial", "error", err)
		return
	}

	smbSession = smbSession

	shareNames, err := smbSession.ListSharenames()
	if err != nil {
		slog.Error("list share names", "error", err)
		return
	}
	fmt.Println(shareNames)
	share, err := smbSession.Mount("labulakalia")
	if err != nil {
		slog.Error("mount dial", "error", err)
		return
	}
	fmt.Println(share)
	share.Umount()
	// smbSession.Logoff()
	slog.Info("umount sussess")
}
