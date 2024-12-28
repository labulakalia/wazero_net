package net

import (
	"context"
	"net"
	"net/http"
)


func init() {
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		var (
			conn *Conn
			err error
		)
		wait := make(chan struct{})
		go func(){
			conn,err = Dial(network, addr)
			wait <- struct{}{}
		}()
		select {
		case <-ctx.Done():
			return nil,ctx.Err()
		case <- wait:
		}
		return conn,err
	}

}
