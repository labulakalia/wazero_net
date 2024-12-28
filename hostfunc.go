package wazero_net

// import (
// 	"context"
// 	"log/slog"
// 	"syscall"
// 	"time"
// 	"wazero_net/errcode"
// 	"wazero_net/util"

// 	"github.com/tetratelabs/wazero/api"
// 	"tinygo.org/x/drivers/netdev"
// )

// func NewExportFunc() *ExportFunc {
// 	return &ExportFunc{
// 		netDev: &NetDev{},
// 	}
// }

// type ExportFunc struct {
// 	netDev netdev.Netdever
// }

// func (e *ExportFunc) Accept(ctx context.Context, m api.Module,
// 	sockfd, r_nfd, r_addrPortPtr, r_addrPortLen uint32) uint32 {
// 	fd, addrPort, err := e.netDev.Accept(int(sockfd))
// 	if err != nil {
// 		slog.Error("accept failed", "err", err)
// 		return errcode.ERR_ACCEPT
// 	}

// 	bytes, err := util.MarshalAddrPort(addrPort)
// 	if err != nil {
// 		slog.Error("marshaladdrport", "err", err)
// 		return errcode.ERR_MARSHAL
// 	}

// 	bytesPtr, err := WriteBytes(m, bytes)
// 	if err != nil {
// 		slog.Error("write bytes failed", "err", err)
// 		return errcode.ERR_WRITE_MEM
// 	}
// 	ok := m.Memory().WriteUint64Le(r_nfd, uint64(fd))
// 	if !ok {
// 		slog.Error("write uint32 failed")
// 		return errcode.ERR_WRITE_MEM
// 	}
// 	ok = m.Memory().WriteUint64Le(r_addrPortPtr, bytesPtr)
// 	if !ok {
// 		slog.Error("write uint32 failed")
// 		return errcode.ERR_WRITE_MEM
// 	}
// 	ok = m.Memory().WriteUint64Le(r_addrPortLen, uint64(len(bytes)))
// 	if !ok {
// 		slog.Error("write uint32 failed")
// 		return errcode.ERR_WRITE_MEM
// 	}
// 	return 0
// }

// func (e *ExportFunc) Addr(ctx context.Context, m api.Module,
// 	r_addrPtr, r_addrLen uint32) uint32 {
// 	return 0
// }

// func (e *ExportFunc) Bind(ctx context.Context, m api.Module,
// 	sockfd, addrPort, addrLen uint32) uint32 {
// 	bytes, err := ReadBytes(m, addrPort, addrLen)
// 	if err != nil {
// 		slog.Error("read uint32 failed", "err", err)
// 		return errcode.ERR_READ_MEM
// 	}
// 	addrport, err := util.UnmarshalAddrPort(bytes)
// 	if err != nil {
// 		slog.Error("unmarshal addr port failed", "err", err)
// 		return errcode.ERR_MARSHAL
// 	}
// 	err = e.netDev.Bind(int(sockfd), addrport)
// 	if err != nil {
// 		slog.Error("bind failed", "err", err)
// 		return errcode.ERR_MARSHAL
// 	}
// 	return 0
// }

// func (e *ExportFunc) Close(ctx context.Context, m api.Module,
// 	sockfd uint32) uint32 {
// 	err := e.netDev.Close(int(sockfd))
// 	if err != nil {
// 		slog.Error("close failed", "err", err)
// 		return errcode.ERR_CLOSE
// 	}
// 	return 0
// }

// func (e *ExportFunc) Connect(ctx context.Context, m api.Module,
// 	sockfd, hostPtr, hostLen, addrPortPtr, addrPortLen uint32) uint32 {
// 	bytes, err := ReadBytes(m, addrPortPtr, addrPortLen)
// 	if err != nil {
// 		slog.Error("mem failed", "err", err)
// 		return errcode.ERR_READ_MEM
// 	}
// 	addrPort, err := util.UnmarshalAddrPort(bytes)
// 	if err != nil {
// 		slog.Error("unmarshal addr port failed", "err", err)
// 		return errcode.ERR_UNMARSHAL
// 	}
// 	hostBytes, err := ReadBytes(m, hostPtr, hostLen)
// 	if err != nil {
// 		slog.Error("mem failed", "err", err)
// 		return errcode.ERR_READ_MEM
// 	}
// 	err = e.netDev.Connect(int(sockfd), util.BytesToString(hostBytes), addrPort)
// 	if err != nil {
// 		slog.Error("connect failed", "err", err)
// 		return errcode.ERR_CONNECT
// 	}
// 	return 0
// }

// func (e *ExportFunc) GetHostByName(ctx context.Context, m api.Module,
// 	namePtr, nameLen, r_addrPtr, r_addrLen uint32) uint32 {
// 	bytes, err := ReadBytes(m, namePtr, nameLen)
// 	if err != nil {
// 		slog.Error("mem failed", "err", err)
// 		return errcode.ERR_READ_MEM
// 	}
// 	addr, err := e.netDev.GetHostByName(util.BytesToString(bytes))
// 	if err != nil {
// 		slog.Error("get host by name failed", "name", string(bytes), "err", err)
// 		return errcode.ERR_GET_HOST_BY_NAME
// 	}
// 	addrBytes, err := addr.MarshalBinary()
// 	if err != nil {
// 		slog.Error("marshal failed", "err", err)
// 		return errcode.ERR_MARSHAL
// 	}
// 	dataPtr, err := WriteBytes(m, addrBytes)
// 	if err != nil {
// 		slog.Error("write bytes failed", "err", err)
// 		return errcode.ERR_MARSHAL
// 	}
// 	ok := m.Memory().WriteUint64Le(r_addrPtr, dataPtr)
// 	if !ok {
// 		slog.Error("write uint32 failed")
// 		return errcode.ERR_WRITE_MEM
// 	}
// 	ok = m.Memory().WriteUint64Le(r_addrLen, uint64(len(addrBytes)))
// 	if !ok {
// 		slog.Error("write uint32 failed")
// 		return errcode.ERR_WRITE_MEM
// 	}
// 	return 0
// }

// func (e *ExportFunc) Listen(ctx context.Context, m api.Module,
// 	sockfd, backlog uint32) uint32 {
// 	err := e.netDev.Listen(int(sockfd), int(backlog))
// 	if err != nil {
// 		slog.Error("listen failed", "err", err)
// 		return errcode.ERR_LISTEN
// 	}
// 	return 0
// }

// func (e *ExportFunc) Recv(ctx context.Context, m api.Module,
// 	sockfd, bufPtr, bufLen uint32, flags, deadline uint64, rn uint32) uint32 {
// 	bufBytes, err := ReadBytes(m, bufPtr, bufLen)
// 	if err != nil {
// 		slog.Error("read bytes failed", "err", err)
// 		return errcode.ERR_READ_MEM
// 	}

// 	n, err := e.netDev.Recv(int(sockfd), bufBytes, int(flags), time.Unix(int64(deadline), 0))
// 	if err != nil {
// 		return errcode.ERR_RECV
// 	}
// 	ok := m.Memory().WriteUint32Le(rn, uint32(n))
// 	if !ok {
// 		return errcode.ERR_WRITE_MEM
// 	}
// 	return 0
// }

// func (e *ExportFunc) Send(ctx context.Context, m api.Module,
// 	sockfd, bufPtr, bufLen uint32, flags, deadline uint64, rn uint32) uint32 {
// 	bufBytes, err := ReadBytes(m, bufPtr, bufLen)
// 	if err != nil {
// 		slog.Error("read bytes failed", "err", err)
// 		return errcode.ERR_READ_MEM
// 	}

// 	n, err := e.netDev.Send(int(sockfd), bufBytes, int(flags), time.Unix(int64(deadline), 0))
// 	if err != nil {
// 		return errcode.ERR_SEND
// 	}
// 	ok := m.Memory().WriteUint64Le(rn, uint64(n))
// 	if !ok {
// 		return errcode.ERR_WRITE_MEM
// 	}
// 	return 0
// }

// func (e *ExportFunc) Socket(ctx context.Context, m api.Module,
// 	domain, stype, protocol, rfd uint32) uint32 {

// 	fd, err := e.netDev.Socket(int(domain), int(stype), int(protocol))
// 	if err != nil {
// 		slog.Error("socket failed", "err", err)
// 		return errcode.ERR_SOCKET
// 	}
// 	ok := m.Memory().WriteUint32Le(rfd, uint32(fd))
// 	if !ok {
// 		return errcode.ERR_WRITE_MEM
// 	}
// 	return 0
// }

// func (e *ExportFunc) SetsockoptByte(ctx context.Context, m api.Module,
// 	sockfd, level, opt uint32, byteData uint32) uint32 {

// 	err := e.netDev.SetSockOpt(int(sockfd), int(level), int(opt), byte(byteData))
// 	if err != nil {
// 		return errcode.ERR_SET_SOCKET_OPT
// 	}
// 	return 0
// }

// func (e *ExportFunc) SetsockoptInt(ctx context.Context, m api.Module,
// 	sockfd, level, opt, intValue uint32) uint32 {
// 	err := e.netDev.SetSockOpt(int(sockfd), int(level), int(opt), int(intValue))
// 	if err != nil {
// 		return errcode.ERR_SET_SOCKET_OPT
// 	}
// 	return 0
// }

// func (e *ExportFunc) SetsockoptInet4Addr(ctx context.Context, m api.Module,
// 	sockfd, level, opt, inet4AddrPtr, inet4AddrLength uint32) uint32 {

// 	inet4Addr, err := ReadBytes(m, inet4AddrPtr, inet4AddrLength)
// 	if err != nil {
// 		return errcode.ERR_READ_MEM
// 	}
// 	if len(inet4Addr) != 4 {
// 		return errcode.ERR_READ_MEM
// 	}

// 	err = e.netDev.SetSockOpt(int(sockfd), int(level), int(opt), *(*[4]byte)(inet4Addr))
// 	if err != nil {
// 		return errcode.ERR_SET_SOCKET_OPT
// 	}
// 	return 0
// }

// func (e *ExportFunc) SetsockoptString(ctx context.Context, m api.Module,
// 	sockfd, level, opt, stringPtr, stringLength uint32) uint32 {

// 	bytes, err := ReadBytes(m, stringPtr, stringLength)
// 	if err != nil {
// 		return errcode.ERR_READ_MEM
// 	}

// 	err = e.netDev.SetSockOpt(int(sockfd), int(level), int(opt), util.BytesToString(bytes))
// 	if err != nil {
// 		return errcode.ERR_SET_SOCKET_OPT
// 	}
// 	return 0
// }

// func (e *ExportFunc) SetsockoptIPMreq(ctx context.Context, m api.Module,
// 	sockfd, level, opt, Multiaddr, Interface uint32) uint32 {

// 	err := e.netDev.SetSockOpt(int(sockfd), int(level), int(opt), &syscall.IPMreq{
// 		Multiaddr: *(*[4]byte)(util.Uint32toByte4(Multiaddr)),
// 		Interface: *(*[4]byte)(util.Uint32toByte4(Interface)),
// 	})
// 	if err != nil {
// 		return errcode.ERR_SET_SOCKET_OPT
// 	}
// 	return 0
// }

// func (e *ExportFunc) SetsockoptIPv6Mreq(ctx context.Context, m api.Module,
// 	sockfd, level, opt, ipv6MreqMultiaddrPtr, ipv6MreqMultiaddrLength, ipv6MreqInterface uint32) uint32 {

// 	ipMreq, err := ReadBytes(m, ipv6MreqMultiaddrPtr, ipv6MreqMultiaddrLength)
// 	if err != nil {
// 		return errcode.ERR_READ_MEM
// 	}
// 	if len(ipMreq) != 16 {
// 		return errcode.ERR_READ_MEM
// 	}

// 	err = e.netDev.SetSockOpt(int(sockfd), int(level), int(opt), &syscall.IPv6Mreq{
// 		Multiaddr: *(*[16]byte)(ipMreq),
// 		Interface: ipv6MreqInterface,
// 	})
// 	if err != nil {
// 		return errcode.ERR_SET_SOCKET_OPT
// 	}
// 	return 0
// }

// func (e *ExportFunc) SetsockoptICMPv6Filter(ctx context.Context, m api.Module,
// 	sockfd, level, opt, icmpv6FilterPtr, icmpv6FilterLength uint32) uint32 {

// 	bytes, err := ReadBytes(m, icmpv6FilterPtr, icmpv6FilterLength)
// 	if err != nil {
// 		return errcode.ERR_READ_MEM
// 	}
// 	if len(bytes) != 32 {
// 		return errcode.ERR_READ_MEM
// 	}
// 	err = e.netDev.SetSockOpt(int(sockfd), int(level), int(opt), &syscall.ICMPv6Filter{
// 		Filt: *(*[8]uint32)(util.BytesToUint32Arr(bytes)),
// 	})
// 	if err != nil {
// 		return errcode.ERR_SET_SOCKET_OPT
// 	}
// 	return 0
// }

// func (e *ExportFunc) SetsockoptLinger(ctx context.Context, m api.Module,
// 	sockfd, level, opt, onoff, linger int32) uint32 {
// 	err := e.netDev.SetSockOpt(int(sockfd), int(level), int(opt), &syscall.Linger{
// 		Onoff:  onoff,
// 		Linger: linger,
// 	})
// 	if err != nil {
// 		return errcode.ERR_SET_SOCKET_OPT
// 	}
// 	return 0
// }

// func (e *ExportFunc) SetsockoptTimeval(ctx context.Context, m api.Module,
// 	sockfd, level, opt, sec uint64, usec uint32) uint32 {
// 	err := e.netDev.SetSockOpt(int(sockfd), int(level), int(opt), &syscall.Timeval{
// 		Sec:       int64(sec),
// 		Usec:      int32(usec),
// 		Pad_cgo_0: [4]byte{0, 0, 0, 0},
// 	})
// 	if err != nil {
// 		return errcode.ERR_SET_SOCKET_OPT
// 	}
// 	return 0
// }
