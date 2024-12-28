package wazero_net

import (
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"syscall"
	"time"

	"tinygo.org/x/drivers/netdev"
)

// TODO support windows

var _ netdev.Netdever = &NetDev{}

type NetDev struct {
}

// Accept implements netdev.Netdever.
func (n *NetDev) Accept(sockfd int) (int, netip.AddrPort, error) {
	slog.Debug("accept", "sockfd", sockfd)
	nfd, addr, err := syscall.Accept(sockfd)
	if err != nil {
		return 0, netip.AddrPort{}, err
	}
	var addrPort netip.AddrPort
	addr4, ok := addr.(*syscall.SockaddrInet4)
	if ok {
		addrPort = netip.AddrPortFrom(netip.AddrFrom4(addr4.Addr), uint16(addr4.Port))
	}
	addr16, ok := addr.(*syscall.SockaddrInet6)
	if ok {
		addrPort = netip.AddrPortFrom(netip.AddrFrom16(addr16.Addr), uint16(addr16.Port))
	}
	return nfd, addrPort, nil
}

// Addr implements netdev.Netdever.
func (n *NetDev) Addr() (netip.Addr, error) {
	return netip.Addr{}, nil

}

// Bind implements netdev.Netdever.
func (n *NetDev) Bind(sockfd int, ip netip.AddrPort) error {
	slog.Debug("accept", "sockfd", sockfd, "ip", ip)
	var sa syscall.Sockaddr
	if ip.Addr().Is4() {
		sa = &syscall.SockaddrInet4{
			Port: int(ip.Port()),  // Listen on this port number
			Addr: ip.Addr().As4(), // Listen to all IPs
		}

	}
	if ip.Addr().Is4() {
		sa = &syscall.SockaddrInet6{
			Port: int(ip.Port()),   // Listen on this port number
			Addr: ip.Addr().As16(), // Listen to all IPs
		}
	}

	return syscall.Bind(sockfd, sa)
}

// Close implements netdev.Netdever.
func (n *NetDev) Close(sockfd int) error {
	slog.Debug("close", "sockfd", sockfd)
	return syscall.Close(sockfd)
}

// Connect implements netdev.Netdever.
func (n *NetDev) Connect(sockfd int, host string, ip netip.AddrPort) error {
	slog.Debug("close", "sockfd", sockfd, "host", host, "ip", ip)
	var sa syscall.Sockaddr
	if ip.Addr().Is4() {
		sa = &syscall.SockaddrInet4{
			Port: int(ip.Port()),  // Listen on this port number
			Addr: ip.Addr().As4(), // Listen to all IPs
		}
	} else {
		sa = &syscall.SockaddrInet6{
			Port: int(ip.Port()),   // Listen on this port number
			Addr: ip.Addr().As16(), // Listen to all IPs
		}
	}
	return syscall.Connect(sockfd, sa)
}

// GetHostByName implements netdev.Netdever.
func (n *NetDev) GetHostByName(name string) (netip.Addr, error) {
	slog.Debug("get host by name", "name", name)
	var addr netip.Addr
	ipAddr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		return addr, err
	}
	if ipAddr.IP.To4() != nil {
		addr = netip.AddrFrom4([4]byte(ipAddr.IP.To4()))

	} else {
		addr = netip.AddrFrom16([16]byte(ipAddr.IP.To16()))
	}
	return addr, nil
}

// Listen implements netdev.Netdever.
func (n *NetDev) Listen(sockfd int, backlog int) error {
	slog.Debug("listen", "sockfd", sockfd, "backlog", backlog)
	return syscall.Listen(sockfd, backlog)
}

// Recv implements netdev.Netdever.
func (n *NetDev) Recv(sockfd int, buf []byte, flags int, deadline time.Time) (int, error) {
	slog.Debug("recv", "sockfd", sockfd, "buf", len(buf), "flags", flags, "deadline", deadline)
	// TODO deadline
	return syscall.Read(sockfd, buf)
}

// Send implements netdev.Netdever.
func (n *NetDev) Send(sockfd int, buf []byte, flags int, deadline time.Time) (int, error) {
	slog.Debug("send", "sockfd", sockfd, "buf", len(buf), "flags", flags, "deadline", deadline)
	// TODO deadline
	return syscall.Write(sockfd, buf)
}

// Socket implements netdev.Netdever.
func (n *NetDev) Socket(domain int, stype int, protocol int) (int, error) {
	slog.Error("send", "domain", domain, "stype", stype, "protocol", protocol)
	return syscall.Socket(domain, stype, protocol)
}

// SetSockOpt implements netdev.Netdever.
func (n *NetDev) SetSockOpt(sockfd int, level int, opt int, value interface{}) error {
	slog.Debug("set sock opt", "level", level, "opt", opt, "value", value)
	switch v := value.(type) {
	case byte:
		return syscall.SetsockoptByte(sockfd, level, opt, v)
	case int:
		return syscall.SetsockoptInt(sockfd, level, opt, v)
	case [4]byte:
		return syscall.SetsockoptInet4Addr(sockfd, level, opt, v)
	case string:
		return syscall.SetsockoptString(sockfd, level, opt, v)
	case *syscall.IPMreq:
		return syscall.SetsockoptIPMreq(sockfd, level, opt, v)
	case *syscall.IPv6Mreq:
		return syscall.SetsockoptIPv6Mreq(sockfd, level, opt, v)
	case *syscall.ICMPv6Filter:
		return syscall.SetsockoptICMPv6Filter(sockfd, level, opt, v)
	case *syscall.Linger:
		return syscall.SetsockoptLinger(sockfd, level, opt, v)
	case *syscall.Timeval:
		return syscall.SetsockoptTimeval(sockfd, level, opt, v)
	}
	return fmt.Errorf("unsupport value type %T", value)
}
