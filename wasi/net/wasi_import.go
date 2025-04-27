//go:build wasip1

package net

//go:wasmimport net conn_dial
func conn_dial(networkPtr, networkLen, address, addressLen, connId uint64) uint64

//go:wasmimport net conn_dial_tls
func conn_dial_tls(networkPtr, networkLen, address, addressLen, connId uint64) uint64

//go:wasmimport net conn_tls_handshake
func conn_tls_handshake(connId uint64) uint64

//go:wasmimport net conn_read
func conn_read(connId, bPtr, bLen, nPtr uint64) uint64

//go:wasmimport net conn_write
func conn_write(connId, bPtr, bLen, nPtr uint64) uint64

//go:wasmimport net conn_close
func conn_close(connId uint64) uint64

// return netip.AddrPort MarshalBinary
//
//go:wasmimport net conn_remote_addr
func conn_remote_addr(connId, addrPtr, addrLenPtr uint64) uint64

//go:wasmimport net conn_local_addr
func conn_local_addr(connId, addrPtr, addrLenPtr uint64) uint64

//go:wasmimport net conn_set_dead_line
func conn_set_dead_line(connId, deadline uint64) uint64

//go:wasmimport net conn_set_read_dead_line
func conn_set_read_dead_line(connId, deadline uint64) uint64

//go:wasmimport net conn_set_write_dead_line
func conn_set_write_dead_line(connId, deadline uint64) uint64

//go:wasmimport net listener_listen
func listener_listen(networkPtr, networkLen, addressPtr, addressLen, listenerIdPtr uint64) uint64

//go:wasmimport net listener_accept
func listener_accept(listenerId, connIdPtr uint64) uint64

//go:wasmimport net listener_close
func listener_close(listenerId uint64) uint64

//go:wasmimport net listener_addr
func listener_addr(listenerId, addrPtr, addrLenPtr uint64) uint64

//go:wasmimport net http_get
func http_get(urlPtr, urlLen uint64) uint64
