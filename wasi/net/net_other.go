//go:build !wasip1

package http

import (
	"crypto/tls"
	"net"
)

func Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}

func DialTLS(network, address string, config *tls.Config) (net.Conn, error) {
	return tls.Dial(network, address, config)
}
