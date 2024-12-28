package net

import (
	"io"
	"log/slog"
	"net"
	"time"
	"wazero_net/errcode"
	"wazero_net/util"
)


func Dial(network,address string) (*Conn, error) {
	slog.Debug("dial", "network",network,"address",address)
	var id uint64
	networkPtr  := util.StringToPtr(&network)
	addressPtr  := util.StringToPtr(&address)
	ret := conn_dial(networkPtr ,uint64(len(network)),
			addressPtr ,uint64(len(address)),
			util.Uint64ToPtr(&id))
	if ret != 0 {
		return nil,util.RetUint64ToError(ret)
	}

	return &Conn{id: uint64(id),network: network},nil
}


var _ net.Conn = &Conn{}
type Conn struct {
	id uint64
	network string
}

func (c *Conn) Read(b []byte) (int, error){
	slog.Debug("conn read", "network",c.network,"id",c.id,"len",len(b))
	var n uint64
	bPtr := util.BytesToPtr(b)
	ret := conn_read(c.id, bPtr,uint64(len(b)), util.Uint64ToPtr(&n))
	if ret == errcode.ERR_CONN_READ_IO_EOF {
		return 0,io.EOF
	}
	err := util.RetUint64ToError(ret)
	if err != nil {
		return 0,err
	}
	return int(n),nil
}


func (c *Conn) Write(b []byte) (int, error){
	slog.Debug("conn write", "network",c.network,"id",c.id,"len",len(b))
	var n uint64
	bPtr := util.BytesToPtr(b)
	err := util.RetUint64ToError(conn_write(c.id, bPtr,uint64(len(b)), util.Uint64ToPtr(&n)))
	if err != nil {
		return 0,err
	}
	return int(n),nil
}

func (c *Conn) Close() error {
	slog.Debug("conn close", "network",c.network,"id",c.id)
	return util.RetUint64ToError(conn_close(c.id))
}

func (c *Conn) RemoteAddr() net.Addr{
	slog.Debug("conn remote addr", "network",c.network,"id",c.id)
	// TODO check data size is enough
	data := memPool.Get().([]byte)
	defer func(){
		data = data[0:0]
		memPool.Put(data)
	}()
	dataPtr := util.BytesToPtr(data)
	dataLength := uint64(len(data))
	err := util.RetUint64ToError(conn_remote_addr(c.id, dataPtr,util.Uint64ToPtr(&dataLength)))
	if err != nil {
		return nil
	}
	data = data[:dataLength]
	var addr net.Addr
	// TODO support more protocol
	switch c.network {
	case "tcp":
		addr,err = net.ResolveTCPAddr(c.network, util.BytesToString(data))
		if err != nil {
			return nil
		}
	case "udp":
		addr,err = net.ResolveUDPAddr(c.network, util.BytesToString(data))
		if err != nil {
			return nil
		}
	}
	return addr
}
func (c *Conn) LocalAddr() net.Addr{
	slog.Debug("conn local addr", "network",c.network,"id",c.id)
	data := memPool.Get().([]byte)
	defer func(){
		data = data[0:0]
		memPool.Put(data)
	}()
	dataPtr := util.BytesToPtr(data)
	dataLen := uint64(len(data))
	err := util.RetUint64ToError(conn_local_addr(c.id, dataPtr,util.Uint64ToPtr(&dataLen)))
	if err != nil {
		return nil
	}
	data = data[:dataLen]
	var addr net.Addr
	switch c.network {
	case "tcp":
		addr,err = net.ResolveTCPAddr(c.network, util.BytesToString(data))
		if err != nil {
			return nil
		}
	case "udp":
		addr,err = net.ResolveUDPAddr(c.network, util.BytesToString(data))
		if err != nil {
			return nil
		}
	}
	return addr
}

func (c *Conn) SetDeadline(t time.Time) error{
	slog.Debug("set dead line addr", "network",c.network,"id",c.id,"time",t)
	return util.RetUint64ToError(conn_set_dead_line(c.id,uint64(t.Unix())))
}
func (c *Conn) SetReadDeadline(t time.Time) error{
	slog.Debug("set read dead line addr", "network",c.network,"id",c.id,"time",t)
	return util.RetUint64ToError(conn_set_read_dead_line(c.id,uint64(t.Unix())))
}
func (c *Conn) SetWriteDeadline(t time.Time) error{
	slog.Debug("set write dead line addr", "network",c.network,"id",c.id,"time",t)
	return util.RetUint64ToError(conn_set_write_dead_line(c.id,uint64(t.Unix())))
}


type Listener struct {
	network string
	id uint64
}

func Listen(network string, address string) (*Listener, error){
	slog.Debug("listen", "network",network,"address",address)
	var id uint64
	networkPtr := util.StringToPtr(&network)
	addressPtr := util.StringToPtr(&address)

	err := util.RetUint64ToError(listener_listen(networkPtr ,uint64(len(network)),
			addressPtr ,uint64(len(address)),
			util.Uint64ToPtr(&id)))
	if err != nil {
		return nil,err
	}
	return &Listener{id: id,network: network},nil
}

// Accept waits for and returns the next connection to the listener.
func(l *Listener)Accept() (net.Conn, error){
	slog.Debug("listen accept", "id",l.id,"network",l.network)
	var connId uint64
	err := util.RetUint64ToError(listener_accept(l.id,util.Uint64ToPtr(&connId)))
	if err != nil {
		return nil,err
	}
	return &Conn{id: connId,network: l.network},nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func(l *Listener)Close() error{
	slog.Debug("listen close", "id",l.id,"network",l.network)
	return util.RetUint64ToError(listener_close(l.id))
}

// Addr returns the listener's network address.
func(l *Listener)Addr() net.Addr {
	slog.Debug("addr", "id",l.id,"network",l.network)
	data := memPool.Get().([]byte)
	defer func(){
		data = data[0:0]
		memPool.Put(data)
	}()
	dataPtr := util.BytesToPtr(data)
	dataLen := uint64(len(data))
	err := util.RetUint64ToError(listener_addr(l.id, dataPtr,util.Uint64ToPtr(&dataLen)))
	if err != nil {
		return nil
	}
	data = data[:dataLen]

	var addr net.Addr
	switch l.network {
	case "tcp":
		addr,err = net.ResolveTCPAddr(l.network, util.BytesToString(data))
		if err != nil {
			return nil
		}
	case "udp":
		addr,err = net.ResolveUDPAddr(l.network, util.BytesToString(data))
		if err != nil {
			return nil
		}
	}
	return addr
}
