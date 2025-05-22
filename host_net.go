package wazero_net

import (
	"context"
	"crypto/tls"

	"errors"
	"fmt"
	"net"

	"sync"
	"time"

	"github.com/labulakalia/wazero_net/util"

	"github.com/tetratelabs/wazero/api"
)

type HostNet struct {
	connLock sync.RWMutex
	connId   uint64
	connMap  map[uint64]net.Conn

	listenerLock sync.RWMutex
	listenerId   uint64
	listenerMap  map[uint64]net.Listener
}

func (h *HostNet) getConn(connId uint64) (net.Conn, error) {

	h.connLock.RLock()
	defer h.connLock.RUnlock()
	conn, ok := h.connMap[connId]
	if !ok {
		return nil, fmt.Errorf("conn id %d not exist ", connId)
	}
	return conn, nil
}

func (h *HostNet) delConn(connId uint64) {
	h.connLock.Lock()
	defer h.connLock.Unlock()
	delete(h.connMap, connId)
}

func (h *HostNet) storeConn(c net.Conn) uint64 {
	h.connLock.Lock()
	defer h.connLock.Unlock()
	h.connId += 1
	h.connMap[h.connId] = c
	return h.connId
}

func (h *HostNet) storeListner(l net.Listener) uint64 {
	h.listenerLock.Lock()
	defer h.listenerLock.Unlock()
	h.listenerId += 1
	h.listenerMap[h.listenerId] = l
	return h.listenerId
}

func (h *HostNet) getListner(listenerId uint64) (net.Listener, error) {
	h.listenerLock.Lock()
	defer h.listenerLock.Unlock()

	listener, ok := h.listenerMap[listenerId]
	if !ok {
		return nil, fmt.Errorf("listener id %d not exist ", listenerId)
	}
	return listener, nil
}

func (h *HostNet) conn_dial(_ context.Context, m api.Module,
	networkPtr, networkLen, addressPtr, addressLen, connIdPtr uint64) uint64 {
	network, err := util.HostReadBytes(m, uint32(networkPtr), uint32(networkLen))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	address, err := util.HostReadBytes(m, uint32(addressPtr), uint32(addressLen))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}

	conn, err := net.Dial(util.BytesToString(network), util.BytesToString(address))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	if util.BytesToString(network) == "tcp" {
		conn.(*net.TCPConn).SetKeepAlive(true)
		conn.(*net.TCPConn).SetKeepAlivePeriod(time.Second * 3)
	}
	newConnId := h.storeConn(conn)

	ok := m.Memory().WriteUint64Le(uint32(connIdPtr), newConnId)
	if !ok {
		return util.HostErrorToUint64(m, errors.New("write connid failed"))
	}
	return 0
}

func (h *HostNet) conn_dial_tls(_ context.Context, m api.Module,
	networkPtr, networkLen, addressPtr, addressLen, connIdPtr uint64) uint64 {
	network, err := util.HostReadBytes(m, uint32(networkPtr), uint32(networkLen))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	address, err := util.HostReadBytes(m, uint32(addressPtr), uint32(addressLen))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}

	conn, err := tls.Dial(util.BytesToString(network), util.BytesToString(address), &tls.Config{
		// InsecureSkipVerify: true,
	})
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	newConnId := h.storeConn(conn)

	ok := m.Memory().WriteUint64Le(uint32(connIdPtr), newConnId)
	if !ok {
		return util.HostErrorToUint64(m, errors.New("write conn id failed"))
	}
	return 0
}

func (h *HostNet) conn_tls_handshake(_ context.Context, m api.Module,
	connId uint64) uint64 {
	conn, err := h.getConn(connId)
	if err != nil {

		return util.HostErrorToUint64(m, err)
	}
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return util.HostErrorToUint64(m, errors.New("tls conn type failed"))
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err = tlsConn.HandshakeContext(ctx)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	return 0
}

func (h *HostNet) conn_read(_ context.Context, m api.Module,
	connId, bPtr, bLen, nPtr uint64) uint64 {
	bytes, err := util.HostReadBytes(m, uint32(bPtr), uint32(bLen))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}

	conn, err := h.getConn(connId)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	err = conn.SetReadDeadline(time.Now().Add(time.Millisecond * 10))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	n, err := conn.Read(bytes)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}

	ok := m.Memory().WriteUint64Le(uint32(nPtr), uint64(n))
	if !ok {
		return util.HostErrorToUint64(m, errors.New("write ptr failed"))
	}
	return 0
}

func (h *HostNet) conn_write(_ context.Context, m api.Module,
	connId, bPtr, bLen, nPtr uint64) uint64 {

	bytes, err := util.HostReadBytes(m, uint32(bPtr), uint32(bLen))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	conn, err := h.getConn(connId)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	err = conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 30))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	n, err := conn.Write(bytes)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	ok := m.Memory().WriteUint64Le(uint32(nPtr), uint64(n))
	if !ok {
		return util.HostErrorToUint64(m, err)
	}
	return 0
}

func (h *HostNet) conn_close(_ context.Context, m api.Module,
	connId uint64) uint64 {
	conn, err := h.getConn(connId)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	h.delConn(connId)
	err = conn.Close()
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	return 0
}

func (h *HostNet) conn_remote_addr(_ context.Context, m api.Module,
	connId, addrPtr, addrLenPtr uint64) uint64 {
	conn, err := h.getConn(connId)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}

	remoteAddr := conn.RemoteAddr().String()

	length, ok := m.Memory().ReadUint64Le(uint32(addrLenPtr))
	if !ok {
		return util.HostErrorToUint64(m, errors.New("read addr len failed"))
	}
	data, err := util.HostReadBytes(m, uint32(addrPtr), uint32(length))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	copy(data, util.StringToBytes(&remoteAddr))
	ok = m.Memory().WriteUint64Le(uint32(addrLenPtr), uint64(len(remoteAddr)))
	if !ok {
		return util.HostErrorToUint64(m, err)
	}
	return 0
}

func (h *HostNet) conn_local_addr(_ context.Context, m api.Module,
	connId, addrPtr, addrLenPtr uint64) uint64 {
	conn, err := h.getConn(connId)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	localAddr := conn.LocalAddr().String()

	length, ok := m.Memory().ReadUint64Le(uint32(addrLenPtr))
	if !ok {
		return util.HostErrorToUint64(m, errors.New("read addr len failed"))
	}
	data, err := util.HostReadBytes(m, uint32(addrPtr), uint32(length))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	copy(data, util.StringToBytes(&localAddr))
	ok = m.Memory().WriteUint64Le(uint32(addrLenPtr), uint64(len(localAddr)))
	if !ok {
		return util.HostErrorToUint64(m, errors.New("write addr len failed"))
	}
	return 0
}

func (h *HostNet) conn_set_dead_line(_ context.Context, m api.Module,
	connId, deadline uint64) uint64 {
	conn, err := h.getConn(connId)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	err = conn.SetDeadline(time.Unix(int64(deadline), 0))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	return 0
}

func (h *HostNet) conn_set_read_dead_line(_ context.Context, m api.Module,
	connId, deadline uint64) uint64 {
	conn, err := h.getConn(connId)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	err = conn.SetReadDeadline(time.Unix(int64(deadline), 0))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	return 0
}

func (h *HostNet) conn_set_write_dead_line(_ context.Context, m api.Module,
	connId, deadline uint64) uint64 {
	conn, err := h.getConn(connId)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	err = conn.SetWriteDeadline(time.Unix(int64(deadline), 0))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	return 0
}

func (h *HostNet) listener_listen(_ context.Context, m api.Module,
	networkPtr, networkLen, addressPtr, addressLen, listenerIdPtr uint64) uint64 {
	network, err := util.HostReadBytes(m, uint32(networkPtr), uint32(networkLen))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	address, err := util.HostReadBytes(m, uint32(addressPtr), uint32(addressLen))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	lis, err := net.Listen(util.BytesToString(network), util.BytesToString(address))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	listenerId := h.storeListner(lis)
	ok := m.Memory().WriteUint64Le(uint32(listenerIdPtr), listenerId)
	if !ok {
		return util.HostErrorToUint64(m, errors.New("write listen id failed"))
	}
	return 0
}

func (h *HostNet) listener_accept(_ context.Context, m api.Module,
	listenerId, connIdPtr uint64) uint64 {
	listener, err := h.getListner(listenerId)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	conn, err := listener.Accept()
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	connId := h.storeConn(conn)
	ok := m.Memory().WriteUint64Le(uint32(connIdPtr), connId)
	if !ok {
		return util.HostErrorToUint64(m, errors.New("write conn id failed"))
	}
	return 0
}

func (h *HostNet) listener_close(_ context.Context, m api.Module,
	listenerId uint64) uint64 {
	listener, err := h.getListner(listenerId)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	err = listener.Close()
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	return 0
}

func (h *HostNet) listener_addr(_ context.Context, m api.Module,
	listenerId, addrPtr, addrLenPtr uint64) uint64 {
	listener, err := h.getListner(listenerId)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	addr := listener.Addr().String()
	ok := m.Memory().Write(uint32(addrPtr), util.StringToBytes(&addr))
	if !ok {
		return util.HostErrorToUint64(m, errors.New("write addr failed"))
	}
	ok = m.Memory().WriteUint64Le(uint32(addrLenPtr), uint64(len(addr)))
	if !ok {
		return util.HostErrorToUint64(m, errors.New("write addr len failed"))
	}
	return 0
}
