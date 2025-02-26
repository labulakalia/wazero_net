package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"

	"github.com/jlaffaye/ftp"
)

func main() {
	ftps()
	// ftps_implicit()
}

func ftps() {
	ftpConn, err := ftp.Dial("127.0.0.1:2121",
		// ftp.DialWithDialFunc(func(network, address string) (net.Conn, error) {
		// 	return net.Dial(network, address)
		// }),
		// ftp.DialWithExplicitTLS(&tls.Config{InsecureSkipVerify: true}),
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
	ftpConn, err := ftp.Dial("127.0.0.1:990",
		ftp.DialWithDialFunc(func(network, address string) (net.Conn, error) {
			conn, err := net.Dial(network, address)
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
	user := "ftpuser"
	password := "admin"
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
