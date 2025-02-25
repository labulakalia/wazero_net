package main

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/jlaffaye/ftp"
)

func main() {
	ftpConn, err := ftp.Dial("127.0.0.1:990",
		// ftp.DialWithDialFunc(func(network, address string) (net.Conn, error) {
		// 	return tls.Dial(network, address, &tls.Config{
		// 		InsecureSkipVerify: true,
		// 	})
		// }),
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
	entries, err := ftpConn.List("/")
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("%+v\n", entries)
}
