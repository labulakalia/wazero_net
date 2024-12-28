package wasm

// func GetBytesFromPtr(ptr uintptr, len uint32) []byte{
// 	res := unsafe.Slice((*byte)(unsafe.Pointer(ptr)), int(len))
// 	return res
// 	// return unsafe.Slice(unsafe.Pointer(&ptr)	,len)
// 	res = make([]byte, 0, uint32(len))
// 	bytes := unsafe.Slice(&ptr, uint32(len))
// 	for i, b := range bytes {
// 		res = append(res, *(*byte)(unsafe.Pointer(b + uintptr(i))))
// 	}
// 	return res
// }

// var _ netdev.Netdever = &NetDev{}

// type NetDev struct {
// }

// // Accept implements netdev.Netdever.
// func (n *NetDev) Accept(sockfd int) (int, netip.AddrPort, error) {

// 	var fd int
// 	var addrPort = netip.AddrPort{}
// 	var addrPortPtr uint32
// 	var addrPortLen uint32

// 	err := errcode.ToError(Accept(uint32(sockfd),
// 		uint32(uintptr(unsafe.Pointer(&fd))),
// 		uint32(uintptr(unsafe.Pointer(&addrPortPtr))),
// 		uint32(uintptr(unsafe.Pointer(&addrPortLen))),
// 	))
// 	if err != nil {
// 		return 0, addrPort, err
// 	}
// 	bytes := GetBytesFromPtr(uintptr(addrPortPtr), addrPortLen)
// 	err = addrPort.UnmarshalBinary(bytes)
// 	if err != nil {
// 		return 0, addrPort, err
// 	}
// 	return fd, addrPort, nil
// }

// // Addr implements netdev.Netdever.
// func (n *NetDev) Addr() (netip.Addr, error) {

// 	var addr = netip.Addr{}
// 	var addrPtr uint32
// 	var addrLen uint32
// 	err := errcode.ToError(Addr(
// 		uint32(uintptr(unsafe.Pointer(&addrPtr))),
// 		uint32(uintptr(unsafe.Pointer(&addrLen))),
// 	))
// 	if err != nil {
// 		return addr, err
// 	}
// 	bytes := GetBytesFromPtr(uintptr(addrPtr), addrLen)
// 	err = addr.UnmarshalBinary(bytes)
// 	if err != nil {
// 		return addr, err
// 	}
// 	return addr, nil
// }

// // Bind implements netdev.Netdever.
// func (n *NetDev) Bind(sockfd int, ip netip.AddrPort) error {
// 	bytes, err := ip.MarshalBinary()
// 	if err != nil {
// 		return err
// 	}
// 	bytesLen := len(bytes)
// 	err = errcode.ToError(Bind(
// 		uint32(sockfd),
// 		uint32(uintptr(unsafe.Pointer(&bytes[0]))),
// 		uint32(bytesLen),
// 	))
// 	return err
// }

// // Close implements netdev.Netdever.
// func (n *NetDev) Close(sockfd int) error {
// 	return errcode.ToError(Close(
// 		uint32(sockfd),
// 	))
// }

// // Connect implements netdev.Netdever.
// func (n *NetDev) Connect(sockfd int, host string, ip netip.AddrPort) error {
// 	ipBytes, err := ip.MarshalBinary()
// 	if err != nil {
// 		return err
// 	}
// 	hostLen := len(host)
// 	ipBytesLen := len(ipBytes)
// 	return errcode.ToError(Connect(uint32(sockfd),
// 		uint32(uintptr(unsafe.Pointer(unsafe.StringData(host)))),
// 		uint32(hostLen),
// 		uint32(uintptr(unsafe.Pointer(&ipBytes[0]))),
// 		uint32(ipBytesLen),
// 	))
// }

// // GetHostByName implements netdev.Netdever.
// func (n *NetDev) GetHostByName(name string) (netip.Addr, error) {

// 	nameLen := len(name)
// 	var addr = netip.Addr{}
// 	var addrPtr uint32
// 	var addrLen uint32
// 	err := errcode.ToError(GetHostByName(
// 		uint32(uintptr(unsafe.Pointer(unsafe.StringData(name)))),
// 		uint32(nameLen),
// 		uint32(uintptr(unsafe.Pointer(&addrPtr))),
// 		uint32(uintptr(unsafe.Pointer(&addrLen))),
// 	))
// 	if err != nil {
// 		return addr, err
// 	}

// 	bytes := GetBytesFromPtr(uintptr(addrPtr), addrLen)
// 	err = addr.UnmarshalBinary(bytes)
// 	if err != nil {
// 		return addr, err
// 	}

// 	os.Stdout.WriteString(fmt.Sprintln("GetHostByName", name, addr))
// 	return addr, nil
// }

// // Listen implements netdev.Netdever.
// func (n *NetDev) Listen(sockfd int, backlog int) error {
// 	return errcode.ToError(Listen(uint32(sockfd), uint32(backlog)))
// }

// // Recv implements netdev.Netdever.
// func (n *NetDev) Recv(sockfd int, buf []byte, flags int, deadline time.Time) (int, error) {

// 	var (
// 		rn  int
// 		err error
// 	)
// 	err = errcode.ToError(Recv(uint32(sockfd),
// 		uint32(uintptr(unsafe.Pointer(&buf[0]))),
// 		uint32(len(buf)),
// 		uint64(flags),
// 		uint64(deadline.Unix()),
// 		uint32(uintptr(unsafe.Pointer(&rn))),
// 	))
// 	if err != nil {
// 		return 0, err
// 	}
// 	return rn, nil
// }

// // Send implements netdev.Netdever.
// func (n *NetDev) Send(sockfd int, buf []byte, flags int, deadline time.Time) (int, error) {
// 	var (
// 		rn  int
// 		err error
// 	)
// 	err = errcode.ToError(Send(uint32(sockfd),
// 		uint32(uintptr(unsafe.Pointer(&buf[0]))),
// 		uint32(len(buf)),
// 		uint64(flags),
// 		uint64(deadline.Unix()),
// 		uint32(uintptr(unsafe.Pointer(&rn))),
// 	))
// 	if err != nil {
// 		return 0, err
// 	}
// 	return rn, nil
// }

// // Socket implements netdev.Netdever.
// func (n *NetDev) Socket(domain int, stype int, protocol int) (int, error) {
// 	var (
// 		rn uint32
// 	)
// 	slog.Error("send", "domain", domain, "stype", stype, "protocol", protocol)
// 	err := errcode.ToError(Socket(uint32(domain), uint32(stype), uint32(protocol), uint32(uintptr(unsafe.Pointer(&rn)))))
// 	return int(rn), err
// }

// func castToBytes[T any](s []T) []byte {
// 	if len(s) == 0 {
// 		return nil
// 	}

// 	size := unsafe.Sizeof(s[0])
// 	return unsafe.Slice((*byte)(unsafe.Pointer(&s[0])), int(size)*len(s))
// }

// // SetSockOpt implements netdev.Netdever.
// func (n *NetDev) SetSockOpt(sockfd int, level int, opt int, value interface{}) error {
// 	switch v := value.(type) {
// 	case byte:
// 		return errcode.ToError(SetsockoptByte(uint32(sockfd), uint32(level), uint32(opt), uint32(v)))
// 	case int:
// 		return errcode.ToError(SetsockoptInt(uint32(sockfd), uint32(level), uint32(opt), uint32(v)))
// 	case [4]byte:
// 		inet4AddrPtr := unsafe.Pointer(&v)
// 		return errcode.ToError(SetsockoptInet4Addr(uint32(sockfd), uint32(level), uint32(opt), uint32(uintptr(inet4AddrPtr)), 4))
// 	case string:
// 		return errcode.ToError(SetsockoptString(uint32(sockfd), uint32(level), uint32(opt), uint32(uintptr(unsafe.Pointer(&v))), uint32(len(v))))
// 		// tinygo unsupport this data
// 		// case *syscall.IPMreq:
// 		// 	return errcode.ToError(SetsockoptIPMreq(sockfd, level, opt, util.Byte4ToUint32(v.Multiaddr[:]), util.Byte4ToUint32(v.Interface[:])))
// 		// case *syscall.IPv6Mreq:
// 		// 	multiAddr := unsafe.Pointer(&v.Multiaddr[:][0])
// 		// 	return errcode.ToError(SetsockoptIPv6Mreq(sockfd, level, opt, uint32(uintptr(multiAddr)), 16, v.Interface))
// 		// case *syscall.ICMPv6Filter:
// 		// 	bytes := []byte{}
// 		// 	for i := 0; i < 8; i++ {
// 		// 		newBytes := util.Uint32toByte4(v.Filt[i])
// 		// 		bytes = append(bytes, newBytes...)
// 		// 	}
// 		// 	return errcode.ToError(SetsockoptICMPv6Filter(sockfd, level, opt, uint32(uintptr(unsafe.Pointer(&bytes[0]))), uint32(len(bytes))))
// 		// case *syscall.Linger:
// 		// 	return errcode.ToError(SetsockoptLinger(sockfd, level, opt, v.Onoff, v.Linger))
// 		// case *syscall.Timeval:
// 		// 	return errcode.ToError(SetsockoptTimeval(sockfd, level, opt, v.Sec, v.Usec))
// 	}
// 	return nil
// }
