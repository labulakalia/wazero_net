package wazero_net

import (
	"net"

	"github.com/tetratelabs/wazero"
)

func InitFuncExport(r wazero.Runtime) wazero.HostModuleBuilder {
	hostNet := &HostNet{
		connMap:     map[uint64]net.Conn{},
		listenerMap: map[uint64]net.Listener{},
		httpResp: map[uint64][]byte{},
	}
	return r.NewHostModuleBuilder("net").
		NewFunctionBuilder().WithFunc(hostNet.conn_dial).Export("conn_dial").
		NewFunctionBuilder().WithFunc(hostNet.conn_dial_tls).Export("conn_dial_tls").
		NewFunctionBuilder().WithFunc(hostNet.conn_tls_handshake).Export("conn_tls_handshake").
		NewFunctionBuilder().WithFunc(hostNet.conn_read).Export("conn_read").
		NewFunctionBuilder().WithFunc(hostNet.conn_write).Export("conn_write").
		NewFunctionBuilder().WithFunc(hostNet.conn_close).Export("conn_close").
		NewFunctionBuilder().WithFunc(hostNet.conn_remote_addr).Export("conn_remote_addr").
		NewFunctionBuilder().WithFunc(hostNet.conn_local_addr).Export("conn_local_addr").
		NewFunctionBuilder().WithFunc(hostNet.conn_set_dead_line).Export("conn_set_dead_line").
		NewFunctionBuilder().WithFunc(hostNet.conn_set_read_dead_line).Export("conn_set_read_dead_line").
		NewFunctionBuilder().WithFunc(hostNet.conn_set_write_dead_line).Export("conn_set_write_dead_line").
		NewFunctionBuilder().WithFunc(hostNet.listener_listen).Export("listener_listen").
		NewFunctionBuilder().WithFunc(hostNet.listener_accept).Export("listener_accept").
		NewFunctionBuilder().WithFunc(hostNet.listener_close).Export("listener_close").
		NewFunctionBuilder().WithFunc(hostNet.listener_addr).Export("listener_addr").
		NewFunctionBuilder().WithFunc(hostNet.round_trip).Export("round_trip").
		NewFunctionBuilder().WithFunc(hostNet.read_resp).Export("read_resp")

}
