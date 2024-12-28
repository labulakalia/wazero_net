package wazero_net

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"
	"wazero_net/errcode"
	"wazero_net/util"

	"github.com/tetratelabs/wazero/api"
)


type HostNet struct {
	connLock sync.RWMutex
	connId uint64
	connMap map[uint64]net.Conn

	listenerLock sync.RWMutex
	listenerId uint64
	listenerMap map[uint64]net.Listener
}

func (h *HostNet)getConn(connId uint64) (net.Conn,error) {
	h.connLock.RLock()
	defer h.connLock.RUnlock()
	conn,ok := h.connMap[connId]
	if !ok {
		return nil,fmt.Errorf("conn id %d not exist ", connId)
	}
	return conn,nil
}

func (h *HostNet)delConn(connId uint64) {
	h.connLock.Lock()
	defer h.connLock.Unlock()
	delete(h.connMap,connId)
}


func (h *HostNet)storeConn(c net.Conn) uint64 {
	h.connLock.Lock()
	defer h.connLock.Unlock()
	h.connId += 1
	h.connMap[h.connId] = c
	return h.connId
}


func (h *HostNet)storeListner(l net.Listener) uint64 {
	h.listenerLock.Lock()
	defer h.listenerLock.Unlock()
	h.listenerId += 1
	h.listenerMap[h.listenerId] = l
	return h.listenerId
}

func (h *HostNet)getListner(listenerId uint64) (net.Listener,error) {
	h.listenerLock.Lock()
	defer h.listenerLock.Unlock()

	listener,ok := h.listenerMap[listenerId]
	if !ok {
		return nil,fmt.Errorf("listener id %d not exist ", listenerId)
	}
	return listener,nil
}

func (h *HostNet)conn_dial(_ context.Context, m api.Module,
	networkPtr,networkLen,addressPtr,addressLen,connIdPtr uint64) uint64 {
	network,err := ReadBytes(m, uint32(networkPtr), uint32(networkLen))
	if err != nil {
		slog.Error("read bytes failed", "err",err)
		return errcode.ERR_READ_MEM
	}
	address,err := ReadBytes(m, uint32(addressPtr), uint32(addressLen))
	if err != nil {
		slog.Error("read bytes failed", "err",err)
		return errcode.ERR_READ_MEM
	}

	conn,err := net.Dial(util.BytesToString(network), util.BytesToString(address))
	if err != nil {
		slog.Error("dial failed","err",err)
		return errcode.ERR_CONN_DIAL
	}
	newConnId := h.storeConn(conn)

	ok := m.Memory().WriteUint64Le(uint32(connIdPtr),newConnId)
	if !ok {
		slog.Error("store conn failed", "newConnId",newConnId)
		return errcode.ERR_WRITE_MEM
	}
	return 0
}


func (h *HostNet)conn_read(_ context.Context, m api.Module,
	connId,bPtr,bLen,nPtr uint64) uint64 {
	bytes,err := ReadBytes(m, uint32(bPtr), uint32(bLen))
	if err != nil {
		slog.Error("read bytes failed", "err",err)
		return errcode.ERR_READ_MEM
	}
	conn,err := h.getConn(connId)
	if err != nil {
		slog.Error("conn not exist failed", "connId",connId)
		return errcode.ERR_CONN_NOT_EXIST
	}
	n,err := conn.Read(bytes)
	if err != nil {
		if errors.Is(err, io.EOF) {
			slog.Debug("read finished")
			return errcode.ERR_CONN_READ_IO_EOF
		}
		slog.Error("read failed", "connId",connId)
		return errcode.ERR_CONN_READ
	}
	ok := m.Memory().WriteUint64Le(uint32(nPtr), uint64(n))
	if !ok {
		return errcode.ERR_WRITE_MEM
	}
	return 0
}


func (h *HostNet)conn_write(_ context.Context, m api.Module,
	connId,bPtr,bLen,nPtr uint64)  uint64 {
	bytes,err := ReadBytes(m, uint32(bPtr), uint32(bLen))
	if err != nil {
		slog.Error("read bytes failed", "err",err)
		return errcode.ERR_READ_MEM
	}
	conn,err := h.getConn(connId)
	if err != nil {
		slog.Error("conn not exist failed", "connId",connId)
		return errcode.ERR_CONN_NOT_EXIST
	}
	n,err := conn.Write(bytes)
	if err != nil {
		slog.Error("write failed", "connId",connId,"err",err)
		return errcode.ERR_CONN_WRITE
	}
	ok := m.Memory().WriteUint64Le(uint32(nPtr), uint64(n))
	if !ok {
		return errcode.ERR_WRITE_MEM
	}
	return 0
}

func (h *HostNet)conn_close(_ context.Context, m api.Module,
	connId uint64) uint64 {
	conn,err := h.getConn(connId)
	if err != nil {
		slog.Error("conn not exist", "connId",connId)
		return errcode.ERR_CONN_NOT_EXIST
	}
	h.delConn(connId)
	err = conn.Close()
	if err != nil {
		slog.Error("close failed", "err",err)
		return errcode.ERR_CONN_CLOSE
	}
	return 0
}

func (h *HostNet)conn_remote_addr(_ context.Context, m api.Module,
	connId ,addrPtr,addrLenPtr uint64) uint64 {
	conn,err := h.getConn(connId)
	if err != nil {
		slog.Error("conn not exist", "connId",connId)
		return errcode.ERR_CONN_NOT_EXIST
	}
	remoteAddr := conn.RemoteAddr().String()
	data,err := ReadBytes(m, uint32(addrPtr), uint32(addrLenPtr))
	if err != nil {
		slog.Error("read bytes failed", "err",err)
		return errcode.ERR_READ_MEM
	}
	copy(data, util.StringToBytes(&remoteAddr))
	ok := m.Memory().WriteUint64Le(uint32(addrLenPtr), uint64(len(remoteAddr)))
	if !ok {
		return errcode.ERR_WRITE_MEM
	}
	return 0
}

func (h *HostNet)conn_local_addr(_ context.Context, m api.Module,
	connId ,addrPtr,addrLenPtr uint64) uint64 {
	conn,err := h.getConn(connId)
	if err != nil {
		slog.Error("conn not exist", "connId",connId)
		return errcode.ERR_CONN_NOT_EXIST
	}
	remoteAddr := conn.LocalAddr().String()
	data,err := ReadBytes(m, uint32(addrPtr), uint32(addrLenPtr))
	if err != nil {
		slog.Error("read bytes failed", "err",err)
		return errcode.ERR_READ_MEM
	}
	copy(data, util.StringToBytes(&remoteAddr))
	ok := m.Memory().WriteUint64Le(uint32(addrLenPtr), uint64(len(remoteAddr)))
	if !ok {
		return errcode.ERR_WRITE_MEM
	}
	return 0
}


func (h *HostNet)conn_set_dead_line(_ context.Context, m api.Module,
	connId ,deadline uint64) uint64 {
	conn,err := h.getConn(connId)
	if err != nil {
		slog.Error("conn not exist", "connId",connId)
		return errcode.ERR_CONN_NOT_EXIST
	}
	err = conn.SetDeadline(time.Unix(int64(deadline), 0))
	if err != nil {
		slog.Error("set dead line failed", "err",err)
		return errcode.ERR_CONN_SET_DEAD_LINE
	}
	return 0
}

func (h *HostNet)conn_set_read_dead_line(_ context.Context, m api.Module,
	connId ,deadline uint64) uint64 {
	conn,err := h.getConn(connId)
	if err != nil {
		slog.Error("conn not exist", "connId",connId)
		return errcode.ERR_CONN_NOT_EXIST
	}
	err = conn.SetReadDeadline(time.Unix(int64(deadline), 0))
	if err != nil {
		slog.Error("set read dead line failed", "err",err)
		return errcode.ERR_CONN_SET_READ_DEAD_LINE
	}
	return 0
}

func (h *HostNet)conn_set_write_dead_line(_ context.Context, m api.Module,
	connId ,deadline uint64) uint64 {
	conn,err := h.getConn(connId)
	if err != nil {
		slog.Error("conn not exist", "connId",connId)
		return errcode.ERR_CONN_NOT_EXIST
	}
	err = conn.SetWriteDeadline(time.Unix(int64(deadline), 0))
	if err != nil {
		slog.Error("set write dead line failed", "err",err)
		return errcode.ERR_CONN_SET_WRITE_DEAD_LINE
	}
	return 0
}


func (h *HostNet)listener_listen(_ context.Context, m api.Module,
	networkPtr,networkLen,addressPtr,addressLen,listenerIdPtr uint64) uint64 {
	network,err := ReadBytes(m, uint32(networkPtr), uint32(networkLen))
	if err != nil {
		return ErrorToRetUint64(m,err)
	}
	address,err := ReadBytes(m, uint32(addressPtr), uint32(addressLen))
	if err != nil {
		return ErrorToRetUint64(m,err)
	}
	lis,err := net.Listen(util.BytesToString(network), util.BytesToString(address))
	if err != nil {
		return ErrorToRetUint64(m,err)
	}
	listenerId := h.storeListner(lis)
	ok := m.Memory().WriteUint64Le(uint32(listenerIdPtr), listenerId)
	if !ok {
		return errcode.ERR_WRITE_MEM
	}
	return 0
}


func (h *HostNet)listener_accept(_ context.Context, m api.Module,
	listenerId,connIdPtr uint64) uint64  {
	listener,err := h.getListner(listenerId)
	if err != nil {
		slog.Error("listener not found", "listenerId",listenerId)
		return errcode.ERR_LISTENER_NOT_EXIST
	}
	conn,err := listener.Accept()
	if err != nil {
		slog.Error("accept failed", "err",err)
		return errcode.ERR_LISTENER_ACCEPT
	}
	connId := h.storeConn(conn)
	ok := m.Memory().WriteUint64Le(uint32(connIdPtr), connId)
	if !ok {
		return errcode.ERR_WRITE_MEM
	}
	return 0
}

func (h *HostNet)listener_close(_ context.Context, m api.Module,
	listenerId uint64) uint64  {
	listener,err := h.getListner(listenerId)
	if err != nil {
		slog.Error("listener not found", "listenerId",listenerId)
		return errcode.ERR_LISTENER_NOT_EXIST
	}
	err = listener.Close()
	if err != nil {
		slog.Error("listener close failed", "err",err)
		return errcode.ERR_LISTENER_CLOSE
	}
	return 0
}

func (h *HostNet)listener_addr(_ context.Context, m api.Module,
	listenerId ,addrPtr,addrLenPtr uint64) uint64  {
	listener,err := h.getListner(listenerId)
	if err != nil {
		slog.Error("listener not found", "listenerId",listenerId)
		return errcode.ERR_LISTENER_NOT_EXIST
	}
	addr := listener.Addr().String()
	ok := m.Memory().Write(uint32(addrPtr), util.StringToBytes(&addr))
	if !ok {
		return errcode.ERR_WRITE_MEM
	}
	ok = m.Memory().WriteUint64Le(uint32(addrLenPtr), uint64(len(addr)))
	if !ok {
		return errcode.ERR_WRITE_MEM
	}
	return 0
}