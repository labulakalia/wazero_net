package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
	wazero_net "github.com/labulakalia/wazero_net/wasi/net"
)

type Conn struct {
	net.Conn
	readTimeout time.Duration
}

func (c *Conn) Read(b []byte) (n int, err error) {
	err = c.Conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return 0, err
	}

	return c.Conn.Read(b)
}

func NewConn(conn net.Conn, readTimeout time.Duration) *Conn {
	return &Conn{
		readTimeout: readTimeout,
		Conn:        conn,
	}
}

func main() {
}

//go:wasmexport ftp_connect
func ftp_connect() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// slog.SetLogLoggerLevel(slog.LevelDebug)
	ftpConn, err := ftp.Dial("127.0.0.1:21",
		ftp.DialWithDialFunc(func(network, address string) (net.Conn, error) {
			conn, err := wazero_net.Dial(network, address)
			if err != nil {
				log.Panicln(err)
			}
			return conn, nil
		}),
		ftp.DialWithDebugOutput(os.Stdout),
		ftp.DialWithExplicitTLS(&tls.Config{InsecureSkipVerify: true}),
	)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("start login")
	user := "ftpuser"
	password := "admin"
	if user == "" && password == "" {
		password = "anonymous"
		user = "anonymous"
	}
	err = ftpConn.Login(user, password)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("start List")
	entries, err := ftpConn.List("/")
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("%+v\n", entries)
}
