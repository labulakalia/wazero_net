package net

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/labulakalia/wazero_net/errcode"
	"github.com/labulakalia/wazero_net/util"
	_ "github.com/labulakalia/wazero_net/wasi/malloc"
)

// TODO Dial add timeout or ctx
func Dial(network, address string) (*Conn, error) {
	slog.Debug("[WASI] dial", "network", network, "address", address)
	var id uint64
	networkPtr := util.StringToPtr(&network)
	addressPtr := util.StringToPtr(&address)
	ret := conn_dial(networkPtr, uint64(len(network)),
		addressPtr, uint64(len(address)),
		util.Uint64ToPtr(&id))
	if ret != 0 {
		return nil, util.RetUint64ToError(ret)
	}

	return &Conn{id: uint64(id), network: network}, nil
}

func DialTls(network, address string) (*Conn, error) {
	slog.Debug("[WASI] dial tls", "network", network, "address", address)
	var id uint64
	networkPtr := util.StringToPtr(&network)
	addressPtr := util.StringToPtr(&address)
	ret := conn_dial_tls(networkPtr, uint64(len(network)),
		addressPtr, uint64(len(address)),
		util.Uint64ToPtr(&id))
	if ret != 0 {
		return nil, util.RetUint64ToError(ret)
	}

	return &Conn{id: uint64(id), network: network}, nil
}

// net/http/transport.go:1714
func (c *Conn) Handshake() error {
	return util.RetUint64ToError(conn_tls_handshake(c.id))
}

// set tc.ConnectionState() tlsState

var _ net.Conn = &Conn{}

type Conn struct {
	*net.TCPConn
	id      uint64
	network string
}

func (c *Conn) Read(b []byte) (int, error) {
	slog.Info("[WASI] conn read", "network", c.network, "id", c.id, "len", len(b))
	var n uint64
	bPtr := util.BytesToPtr(b)
	ret := conn_read(c.id, bPtr, uint64(len(b)), util.Uint64ToPtr(&n))
	if ret == errcode.ERR_CONN_READ_IO_EOF {
		return int(n), io.EOF
	}
	err := util.RetUint64ToError(ret)
	if err != nil {
		return 0, err
	}
	slog.Info("read success","n",n)
	time.Sleep(time.Millisecond)
	return int(n), nil
}

func (c *Conn) Write(b []byte) (int, error) {
	slog.Info("[WASI] conn write", "network", c.network, "id", c.id, "len", len(b))
	var n uint64
	bPtr := util.BytesToPtr(b)
	err := util.RetUint64ToError(conn_write(c.id, bPtr, uint64(len(b)), util.Uint64ToPtr(&n)))
	if err != nil {
		return 0, err
	}
	time.Sleep(time.Millisecond)
	return int(n), nil
}

func (c *Conn) Close() error {
	slog.Debug("[WASI] conn close", "network", c.network, "id", c.id)
	return util.RetUint64ToError(conn_close(c.id))
}

func (c *Conn) RemoteAddr() net.Addr {
	slog.Info("[WASI] conn remote addr", "network", c.network, "id", c.id)
	// TODO check data size is enough
	data := util.MemPool.Get().([]byte)
	defer func() {

		util.MemPool.Put(data)
	}()
	dataPtr := util.BytesToPtr(data)
	dataLength := uint64(len(data))
	err := util.RetUint64ToError(conn_remote_addr(c.id, dataPtr, util.Uint64ToPtr(&dataLength)))
	if err != nil {
		return nil
	}

	var addr net.Addr
	// TODO support more protocol
	switch c.network {
	case "tcp":
		addr, err = net.ResolveTCPAddr(c.network, util.BytesToString(data[:dataLength]))
		if err != nil {
			return nil
		}
	case "udp":
		addr, err = net.ResolveUDPAddr(c.network, util.BytesToString(data[:dataLength]))
		if err != nil {
			return nil
		}
	}
	return addr
}
func (c *Conn) LocalAddr() net.Addr {
	slog.Info("[WASI] conn local addr", "network", c.network, "id", c.id)
	data := util.MemPool.Get().([]byte)
	defer func() {
		util.MemPool.Put(data)
	}()
	fmt.Println("data", len(data))
	dataPtr := util.BytesToPtr(data)
	dataLen := uint64(len(data))
	err := util.RetUint64ToError(conn_local_addr(c.id, dataPtr, util.Uint64ToPtr(&dataLen)))
	if err != nil {
		slog.Error("read local addr failed","err",err)
		return nil
	}

	var addr net.Addr
	switch c.network {
	case "tcp":
		addr, err = net.ResolveTCPAddr(c.network, util.BytesToString(data[:dataLen]))
		if err != nil {
			return nil
		}
	case "udp":
		addr, err = net.ResolveUDPAddr(c.network, util.BytesToString(data[:dataLen]))
		if err != nil {
			return nil
		}
	}
	return addr
}

func (c *Conn) SetDeadline(t time.Time) error {
	slog.Debug("[WASI] set dead line addr", "network", c.network, "id", c.id, "time", t)
	return util.RetUint64ToError(conn_set_dead_line(c.id, uint64(t.Unix())))
}
func (c *Conn) SetReadDeadline(t time.Time) error {
	slog.Debug("[WASI] set read dead line addr", "network", c.network, "id", c.id, "time", t)
	return util.RetUint64ToError(conn_set_read_dead_line(c.id, uint64(t.Unix())))
}
func (c *Conn) SetWriteDeadline(t time.Time) error {
	slog.Debug("[WASI] set write dead line addr", "network", c.network, "id", c.id, "time", t)
	return util.RetUint64ToError(conn_set_write_dead_line(c.id, uint64(t.Unix())))
}

type Listener struct {
	network string
	id      uint64
}

func Listen(network string, address string) (*Listener, error) {
	slog.Debug("[WASI] listen", "network", network, "address", address)
	var id uint64
	networkPtr := util.StringToPtr(&network)
	addressPtr := util.StringToPtr(&address)

	err := util.RetUint64ToError(listener_listen(networkPtr, uint64(len(network)),
		addressPtr, uint64(len(address)),
		util.Uint64ToPtr(&id)))
	if err != nil {
		return nil, err
	}
	return &Listener{id: id, network: network}, nil
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (net.Conn, error) {
	slog.Debug("[WASI] listen accept", "id", l.id, "network", l.network)
	var connId uint64
	err := util.RetUint64ToError(listener_accept(l.id, util.Uint64ToPtr(&connId)))
	if err != nil {
		return nil, err
	}
	return &Conn{id: connId, network: l.network}, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *Listener) Close() error {
	slog.Debug("[WASI] listen close", "id", l.id, "network", l.network)
	return util.RetUint64ToError(listener_close(l.id))
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	slog.Debug("[WASI] addr", "id", l.id, "network", l.network)
	data := util.MemPool.Get().([]byte)
	defer func() {
		data = data[0:0]
		util.MemPool.Put(data)
	}()
	dataPtr := util.BytesToPtr(data)
	dataLen := uint64(len(data))
	err := util.RetUint64ToError(listener_addr(l.id, dataPtr, util.Uint64ToPtr(&dataLen)))
	if err != nil {
		return nil
	}
	data = data[:dataLen]

	var addr net.Addr
	switch l.network {
	case "tcp":
		addr, err = net.ResolveTCPAddr(l.network, util.BytesToString(data))
		if err != nil {
			return nil
		}
	case "udp":
		addr, err = net.ResolveUDPAddr(l.network, util.BytesToString(data))
		if err != nil {
			return nil
		}
	}
	return addr
}

func HttpGet(url string) {
	urlPtr := util.StringToPtr(&url)
	http_get(urlPtr, uint64(len(url)))
}
