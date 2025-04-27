package main

import (
	"fmt"
	"log/slog"

	_ "github.com/labulakalia/wazero_net/wasi/malloc"
	wasi_net "github.com/labulakalia/wazero_net/wasi/net"
	"github.com/medianexapp/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {}

var sftpClient *sftp.Client

//go:wasmexport sftp_connect
func _sftp_connect() {
	// slog.SetLogLoggerLevel(slog.LevelDebug)
	if sftpClient != nil {
		fmt.Println("old client", sftpClient)
		files, err := sftpClient.ReadDir("/")
		if err != nil {
			slog.Error("read dir failed", "err", err)
			return
		}
		fmt.Println("files", len(files))
		return
	}
	conn, err := wasi_net.Dial("tcp", "127.0.0.1:22")
	if err != nil {
		slog.Error("dial failed", "err", err)
		return
	}

	_, _, _, err = ssh.NewClientConn(conn, "127.0.0.1:22", &ssh.ClientConfig{
		User:            "labulakalia",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password("109097"),
			// ssh.PublicKeys(signer),
		},
		// Timeout: time.Second * 3,
	})
	if err != nil {
		slog.Error("client conn failed", "err", err)
		fmt.Println("client conn failed", err)
		return
	}

	// sshClient := ssh.NewClient(c, chans, reqs)
	// client, err := sftp.NewClient(sshClient)
	// if err != nil {
	// 	slog.Error("new sftp client failed", "err", err)
	// 	fmt.Println("new sftp client failed", err)
	// 	return
	// }

	// sftpClient = client

	// files, err := sftpClient.ReadDir("/")
	// if err != nil {
	// 	slog.Error("read dir failed", "err", err)
	// 	return
	// }
	// fmt.Println("files", len(files))

}
