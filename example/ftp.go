package main

import (
	"log"
	"net"
	"time"

	"github.com/jlaffaye/ftp"
	wasi_net "github.com/labulakalia/wazero_net/wasi/net"
)

func main() {}

//go:wasmexport ftp_connect
func ftp_connect() {
	ftpConn, err := ftp.Dial("127.0.0.1:2121", ftp.DialWithDialFunc(func(network, address string) (net.Conn, error) {
		return wasi_net.Dial(network, address)
	}),
		ftp.DialWithTimeout(time.Second*3),
	)
	if err != nil {
		log.Panicln(err)
	}
	user := "user"
	password := "passwd"
	if user == "" && password == "" {
		password = "anonymous"
		user = "anonymous"
	}
	err = ftpConn.Login(user, password)
	if err != nil {
		log.Panicln(err)
	}
	entries, err := ftpConn.List("/")
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("%+v\n", entries)
}
