//go:build wasip1

package net

import (
	"crypto/tls"
	"io"
	"log/slog"
	"net"
	"net/netip"
	"runtime"
	"strings"
	"time"

	"github.com/labulakalia/wazero_net/util"
)

// TODO Dial add timeout or ctx
func Dial(network, address string) (net.Conn, error) {
	slog.Debug("[WASI] dial", "network", network, "address", address)
	var id uint64
	networkPtr := util.StringToPtr(&network)
	addressPtr := util.StringToPtr(&address)
	time.Sleep(0)
	ret := conn_dial(networkPtr, uint64(len(network)),
		addressPtr, uint64(len(address)),
		util.Uint64ToPtr(&id))
	time.Sleep(0)
	if ret != 0 {
		return nil, util.RetUint64ToError(ret)
	}

	return &Conn{id: uint64(id), network: network}, nil
}

func DialTLS(network, address string, _ *tls.Config) (*Conn, error) {
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

	return &Conn{id: id, network: network}, nil
}

var _ net.Conn = &Conn{}

type Conn struct {
	id      uint64
	network string

	readDeadline  time.Time
	writeDeadline time.Time
}

func (c *Conn) Read(b []byte) (int, error) {
	slog.Debug("[WASI] conn read", "network", c.network, "id", c.id, "len", len(b))
	defer func() {
		if !c.readDeadline.IsZero() {
			c.readDeadline = time.Time{}
		}
	}()
	var n uint64
	bPtr := util.BytesToPtr(b)
reply:
	runtime.Gosched()
	err := util.RetUint64ToError(conn_read(c.id, bPtr, uint64(len(b)), util.Uint64ToPtr(&n)))
	runtime.Gosched()
	if err != nil {
		if err.Error() == "EOF" {
			return int(n), io.EOF
		} else {
			if strings.Contains(err.Error(), "i/o timeout") {
				if !c.readDeadline.IsZero() && time.Now().After(c.readDeadline) {
					slog.Error("read errzero", "err", err)
					return 0, err
				}
				slog.Debug("read err retry", "err", err)
				goto reply
			}
			return 0, err
		}
	}
	slog.Debug("read success", "count", n, "remoteAddr", c.RemoteAddr())
	return int(n), nil
}

func (c *Conn) Write(b []byte) (int, error) {
	slog.Debug("[WASI] conn write", "network", c.network, "id", c.id, "len", len(b))
	defer func() {
		if !c.readDeadline.IsZero() {
			c.readDeadline = time.Time{}
		}
	}()
	var n uint64
	bPtr := util.BytesToPtr(b)
reply:
	err := util.RetUint64ToError(conn_write(c.id, bPtr, uint64(len(b)), util.Uint64ToPtr(&n)))
	runtime.Gosched()
	if err != nil {
		if strings.Contains(err.Error(), "i/o timeout") {
			if !c.writeDeadline.IsZero() && time.Now().After(c.writeDeadline) {
				return 0, err
			}
			goto reply
		}
		return 0, err
	}
	slog.Debug("read success", "n", n)

	return int(n), nil
}

func (c *Conn) Close() error {
	slog.Debug("[WASI] conn close", "network", c.network, "id", c.id)

	return util.RetUint64ToError(conn_close(c.id))
}

func (c *Conn) RemoteAddr() net.Addr {
	slog.Debug("[WASI] conn remote addr", "network", c.network, "id", c.id)
	// TODO check data size is enough
	data := util.MemPool.Get().([]byte)
	defer func() {
		util.MemPool.Put(data)
	}()

	dataPtr := util.BytesToPtr(data)
	dataLength := uint64(len(data))
	err := util.RetUint64ToError(conn_remote_addr(c.id, dataPtr, util.Uint64ToPtr(&dataLength)))
	if err != nil {
		slog.Error("remote addr failed", "err", err)
		return nil
	}
	addrPortStr := util.BytesToString(data[:dataLength])
	addrPort, err := netip.ParseAddrPort(addrPortStr)
	if err != nil {
		slog.Error("parse failed", "err", err)
		return nil
	}
	var addr net.Addr
	// TODO support more protocol
	switch c.network {
	case "tcp":
		addr = &net.TCPAddr{
			IP:   addrPort.Addr().AsSlice(),
			Zone: addrPort.Addr().Zone(),
			Port: int(addrPort.Port()),
		}
	case "udp":
		addr = &net.UDPAddr{
			IP:   addrPort.Addr().AsSlice(),
			Zone: addrPort.Addr().Zone(),
			Port: int(addrPort.Port()),
		}
	}
	return addr
}
func (c *Conn) LocalAddr() net.Addr {
	slog.Debug("[WASI] conn local addr", "network", c.network, "id", c.id)
	data := util.MemPool.Get().([]byte)
	defer func() {
		util.MemPool.Put(data)
	}()

	dataPtr := util.BytesToPtr(data)
	dataLen := uint64(len(data))
	err := util.RetUint64ToError(conn_local_addr(c.id, dataPtr, util.Uint64ToPtr(&dataLen)))
	if err != nil {
		slog.Error("read local addr failed", "err", err)
		return nil
	}
	addrPortStr := util.BytesToString(data[:dataLen])
	addrPort, err := netip.ParseAddrPort(addrPortStr)
	if err != nil {
		slog.Error("parse failed", "err", err)
		return nil
	}
	var addr net.Addr
	switch c.network {
	case "tcp":
		addr = &net.TCPAddr{
			IP:   addrPort.Addr().AsSlice(),
			Zone: addrPort.Addr().Zone(),
			Port: int(addrPort.Port()),
		}
	case "udp":
		addr = &net.UDPAddr{
			IP:   addrPort.Addr().AsSlice(),
			Zone: addrPort.Addr().Zone(),
			Port: int(addrPort.Port()),
		}
	}
	return addr
}

func (c *Conn) SetDeadline(t time.Time) error {
	slog.Debug("[WASI] set dead line addr", "network", c.network, "id", c.id, "time", t)
	c.readDeadline = t
	c.writeDeadline = t
	return nil
}
func (c *Conn) SetReadDeadline(t time.Time) error {
	slog.Debug("[WASI] set read dead line addr", "network", c.network, "id", c.id, "time", t)
	c.readDeadline = t
	return nil
}
func (c *Conn) SetWriteDeadline(t time.Time) error {
	slog.Debug("[WASI] set write dead line addr", "network", c.network, "id", c.id, "time", t)
	c.readDeadline = t
	return nil
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
