package main

import (
	"context"
	"log/slog"

	_ "github.com/labulakalia/wazero_net/wasi/malloc"
	wasi_net "github.com/labulakalia/wazero_net/wasi/net"
	"github.com/medianexapp/go-smb2"
)

func main() {}

var session *smb2.Session

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

	slog.Info("conn success")
	smbSession, err := smbDialer.DialConn(context.Background(), conn, addr)
	if err != nil {
		slog.Error("failed to dial", "error", err)
		return
	}
	session = smbSession
	slog.Info("dial conn", "smbSession", session)
	shareNames, err := session.ListSharenames()
	if err != nil {
		slog.Error("failed to list share names", "error", err)
		return
	}
	slog.Info("msg string", "names", shareNames)
	share, err := session.Mount("labulakalia")
	if err != nil {
		slog.Error("failed to list share names", "error", err)
		return
	}
	dirs, err := share.ReadDir(".")
	if err != nil {
		slog.Error("failed to read dir", "error", err)
		return
	}
	for _, dir := range dirs {
		slog.Info("dir", "name", dir.Name())
	}
	slog.Info("smb exit")

}
