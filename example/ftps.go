package main

import (
	"crypto/tls"
	"log"

	"github.com/jlaffaye/ftp"
)

func main() {
	ftpConn, err := ftp.Dial("127.0.0.1:990",
		ftp.DialWithTLS(&tls.Config{InsecureSkipVerify: true}),
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
