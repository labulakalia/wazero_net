package net

//go:wasmimport net conn_dial
//go:noescape
func conn_dial(networkPtr, networkLen, address, addressLen, connId uint64) uint64

//go:wasmimport net conn_dial_tls
//go:noescape
func conn_dial_tls(networkPtr, networkLen, address, addressLen, connId uint64) uint64

//go:wasmimport net conn_tls_handshake
//go:noescape
func conn_tls_handshake(connId uint64) uint64

//go:wasmimport net conn_read
//go:noescape
func conn_read(connId, bPtr, bLen, nPtr uint64) uint64

//go:wasmimport net conn_write
//go:noescape
func conn_write(connId, bPtr, bLen, nPtr uint64) uint64

//go:wasmimport net conn_close
//go:noescape
func conn_close(connId uint64) uint64

// return netip.AddrPort MarshalBinary
//
//go:wasmimport net conn_remote_addr
//go:noescape
func conn_remote_addr(connId, addrPtr, addrLenPtr uint64) uint64

//go:wasmimport net conn_local_addr
//go:noescape
func conn_local_addr(connId, addrPtr, addrLenPtr uint64) uint64

//go:wasmimport net conn_set_dead_line
//go:noescape
func conn_set_dead_line(connId, deadline uint64) uint64

//go:wasmimport net conn_set_read_dead_line
//go:noescape
func conn_set_read_dead_line(connId, deadline uint64) uint64

//go:wasmimport net conn_set_write_dead_line
//go:noescape
func conn_set_write_dead_line(connId, deadline uint64) uint64

//go:wasmimport net listener_listen
//go:noescape
func listener_listen(networkPtr, networkLen, addressPtr, addressLen, listenerIdPtr uint64) uint64

//go:wasmimport net listener_accept
//go:noescape
func listener_accept(listenerId, connIdPtr uint64) uint64

//go:wasmimport net listener_close
//go:noescape
func listener_close(listenerId uint64) uint64

//go:wasmimport net listener_addr
//go:noescape
func listener_addr(listenerId, addrPtr, addrLenPtr uint64) uint64

//go:wasmimport net http_get
//go:noescape
func http_get(urlPtr, urlLen uint64) uint64
