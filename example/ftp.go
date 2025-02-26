package main

import (
	"crypto/tls"
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
	ftp_conn()
	ftps()
	ftps_implicit()
}

func ftp_conn() {
	slog.Info("ftp connect")
	ftpConn, err := ftp.Dial("127.0.0.1:2121",
		ftp.DialWithDialFunc(func(network, address string) (net.Conn, error) {
			return wazero_net.Dial(network, address)
		}),
		ftp.DialWithDebugOutput(os.Stdout),
	)

	if err != nil {
		log.Panicln(err)
	}
	user := "user"
	password := "passwd"
	err = ftpConn.Login(user, password)
	if err != nil {
		log.Panicln(err)
	}
	entries, err := ftpConn.List("/")
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("ftps %+v\n", entries)
}

func ftps() {
	slog.Info("ftps connect")
	ftpConn, err := ftp.Dial("127.0.0.1:21",
		ftp.DialWithDialFunc(func(network, address string) (net.Conn, error) {
			return wazero_net.Dial(network, address)
		}),
		ftp.DialWithExplicitTLS(&tls.Config{InsecureSkipVerify: true}),
		ftp.DialWithDebugOutput(os.Stdout),
	)

	if err != nil {
		log.Panicln(err)
	}
	user := "user"
	password := "passwd"
	err = ftpConn.Login(user, password)
	if err != nil {
		log.Panicln(err)
	}
	entries, err := ftpConn.List("/")
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("ftps %+v\n", entries)
}

func ftps_implicit() {
	slog.Info("ftps_implicit connect")
	ftpConn, err := ftp.Dial("127.0.0.1:990",
		ftp.DialWithDialFunc(func(network, address string) (net.Conn, error) {
			conn, err := wazero_net.Dial(network, address)
			if err != nil {
				return nil, err
			}
			return tls.Client(conn, &tls.Config{InsecureSkipVerify: true}), nil
		}),
		ftp.DialWithTLS(&tls.Config{InsecureSkipVerify: true}),
		ftp.DialWithDebugOutput(os.Stdout),
	)

	if err != nil {
		log.Panicln(err)
	}
	user := "user"
	password := "passwd"
	err = ftpConn.Login(user, password)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("ftp conn list")
	entries, err := ftpConn.List("/")
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("ftps_implicit %+v\n", entries)
}
