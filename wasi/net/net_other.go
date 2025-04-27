//go:build !wasip1

package http

import "net"

func Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}
