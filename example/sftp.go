package main

import (
	"fmt"
	"log/slog"
	"time"

	wasi_net "github.com/labulakalia/wazero_net/wasi/net"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {}

//go:wasmexport sftp_connect
func sftp_connect() {
	conn, err := wasi_net.Dial("tcp", "127.0.0.1:22")
	if err != nil {
		slog.Error("dial failed", "err", err)
		return
	}
	defer conn.Close()

	c, chans, reqs, err := ssh.NewClientConn(conn, "127.0.0.1:22", &ssh.ClientConfig{
		User:            "labulakalia",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password("109097"),
			// ssh.PublicKeys(signer),
		},
		Timeout: time.Second * 3,
	})
	if err != nil {
		slog.Error("client conn failed", "err", err)
		return
	}
	sshClient := ssh.NewClient(c, chans, reqs)
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		slog.Error("new sftp client failed", "err", err)
		return
	}
	files, err := client.ReadDir("/etc")
	if err != nil {
		slog.Error("read dir failed", "err", err)
		return
	}
	for _, file := range files {
		fmt.Printf("file: %v\n", file.Name())
	}
}
