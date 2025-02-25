package main

import (
	"crypto/tls"
	"log"
	"net"
	"time"

	"github.com/jlaffaye/ftp"
	wasi_net "github.com/labulakalia/wazero_net/wasi/net"
	// wasi_net "github.com/labulakalia/wazero_net/wasi/net"
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
	ftpConn, err := ftp.Dial("127.0.0.1:21", ftp.DialWithDialFunc(func(network, address string) (net.Conn, error) {
		conn, err := wasi_net.Dial(network, address)
		if err != nil {
			return nil, err
		}
		return conn,err
		// return NewConn(conn, time.Second), nil
	}),
		ftp.DialWithExplicitTLS(&tls.Config{InsecureSkipVerify: true}),
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
