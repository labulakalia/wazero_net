package wasm

//go:wasmimport net Accept
//go:noescape
func Accept(sockfd, r_nfd, r_addrPortPtr, r_addrPortLen uint32) uint32

//go:wasmimport net Addr
//go:noescape
func Addr(r_addrPtr, r_addrLen uint32) uint32

//go:wasmimport net Bind
//go:noescape
func Bind(sockfd, addrPort, addrLen uint32) uint32

//go:wasmimport net Close
//go:noescape
func Close(sockfd uint32) uint32

//go:wasmimport net Connect
//go:noescape
func Connect(sockfd, hostPtr, hostLen, addrPortPtr, addrPortLen uint32) uint32

//go:wasmimport net GetHostByName
//go:noescape
func GetHostByName(namePtr, nameLen, r_addrPtr, r_addrLen uint32) uint32

//go:wasmimport net Listen
//go:noescape
func Listen(sockfd, backlog uint32) uint32

//go:wasmimport net Recv
//go:noescape
func Recv(sockfd, bufPtr, bufLen uint32, flags, deadline uint64, rn uint32) uint32

//go:wasmimport net Send
//go:noescape
func Send(sockfd, bufPtr, bufLen uint32, flags, deadline uint64, rn uint32) uint32

//go:wasmimport net Socket
//go:noescape
func Socket(domain, stype, protocol, rfd uint32) uint32

//go:wasmimport net SetsockoptByte
//go:noescape
func SetsockoptByte(sockfd, level, opt uint32, byteData uint32) uint32

//go:wasmimport net SetsockoptInt
//go:noescape
func SetsockoptInt(sockfd, level, opt, intValue uint32) uint32

//go:wasmimport net SetsockoptInet4Addr
//go:noescape
func SetsockoptInet4Addr(sockfd, level, opt, inet4AddrPtr, inet4AddrLength uint32) uint32

//go:wasmimport net SetsockoptString
//go:noescape
func SetsockoptString(sockfd, level, opt, stringPtr, stringLength uint32) uint32

//go:wasmimport net SetsockoptIPMreq
//go:noescape
func SetsockoptIPMreq(sockfd, level, opt, Multiaddr, Interface uint32) uint32

//go:wasmimport net SetsockoptIPv6Mreq
//go:noescape
func SetsockoptIPv6Mreq(sockfd, level, opt, ipv6MreqMultiaddrPtr, ipv6MreqMultiaddrLength, ipv6MreqInterface uint32) uint32

//go:wasmimport net SetsockoptICMPv6Filter
//go:noescape
func SetsockoptICMPv6Filter(sockfd, level, opt, icmpv6FilterPtr, icmpv6FilterLength uint32) uint32

//go:wasmimport net SetsockoptLinger
//go:noescape
func SetsockoptLinger(sockfd, level, opt, onoff, linger uint32) uint32

//go:wasmimport net SetsockoptTimeval
//go:noescape
func SetsockoptTimeval(sockfd, level, opt, sec uint64, usec uint32) uint32
