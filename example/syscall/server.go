package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type netSocket struct {
	fd int
}

func (ns netSocket) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	n, err := syscall.Read(ns.fd, p)
	if err != nil {
		n = 0
	}
	return n, err
}
func (ns netSocket) Write(p []byte) (int, error) {
	n, err := syscall.Write(ns.fd, p)
	if err != nil {
		n = 0
	}
	return n, err
}
func (ns *netSocket) Accept() (*netSocket, error) {
	nfd, _, err := syscall.Accept(ns.fd)
	if err == nil {
		syscall.CloseOnExec(nfd)
	}
	if err != nil {
		return nil, err
	}
	return &netSocket{nfd}, nil
}
func (ns *netSocket) Close() error {
	return syscall.Close(ns.fd)
}

func newNetSocket(ip net.IP, port int) (*netSocket, error) {
	// ForkLock 文档指明需要加锁
	syscall.ForkLock.Lock()
	// 这里第一个参数我们使用syscall.AF_INET, IPv4的地址族。
	// 第二个参数指明是数据流方式，也就是TCP的方式。
	// 第三个参数使用SOCK_STREAM默认协议。
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, os.NewSyscallError("socket", err)
	}
	syscall.ForkLock.Unlock()
	// 建立了Socket，并且得到了文件描述符，我们可以设置一些选项，
	// 比如可重用的地址
	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		syscall.Close(fd)
		return nil, os.NewSyscallError("setsockopt", err)
	}
	// 绑定指定的地址和端口
	sa := &syscall.SockaddrInet4{Port: port}
	copy(sa.Addr[:], ip)
	if err = syscall.Bind(fd, sa); err != nil {
		return nil, os.NewSyscallError("bind", err)
	}
	// 开始监听客户端的连接请求
	if err = syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		return nil, os.NewSyscallError("listen", err)
	}
	return &netSocket{fd: fd}, nil
}

func main() {
	ipFlag := flag.String("ip_addr", "127.0.0.1", "监听的地址")
	portFlag := flag.Int("port", 8080, "监听的端口")
	flag.Parse()
	ip := net.ParseIP(*ipFlag)
	socket, err := newNetSocket(ip, *portFlag)
	if err != nil {
		panic(err)
	}
	defer socket.Close()
	log.Printf("http addr: http://%s:%d", ip, *portFlag)
	for {
		// 开始等待客户端的连接
		rw, e := socket.Accept()
		log.Printf("incoming connection")
		if e != nil {
			panic(e)
		}
		log.Println("accept")
		// Read request
		log.Print("reading request")
		req, err := parseRequest(rw)
		log.Print("request: ", req)
		if err != nil {
			panic(err)
		}
		// Write response
		log.Print("writing response")
		io.WriteString(rw, "HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/html; charset=utf-8\r\n"+
			"Content-Length: 20\r\n"+
			"\r\n"+
			"<h1>hello world</h1>")
		if err != nil {
			log.Print(err.Error())
			continue
		}
	}
}

type request struct {
	method string // GET, POST, etc.
	header textproto.MIMEHeader
	body   []byte
	uri    string // The raw URI from the request
	proto  string // "HTTP/1.1"
}

func parseRequest(c *netSocket) (*request, error) {
	b := bufio.NewReader(*c)
	tp := textproto.NewReader(b)
	req := new(request)
	// First line: parse "GET /index.html HTTP/1.0"
	var s string
	s, _ = tp.ReadLine()
	sp := strings.Split(s, " ")
	req.method, req.uri, req.proto = sp[0], sp[1], sp[2]
	// Parse headers
	mimeHeader, _ := tp.ReadMIMEHeader()
	req.header = mimeHeader
	// Parse body
	if req.method == "GET" || req.method == "HEAD" {
		return req, nil
	}
	if len(req.header["Content-Length"]) == 0 {
		return nil, errors.New("no content length")
	}
	length, err := strconv.Atoi(req.header["Content-Length"][0])
	if err != nil {
		return nil, err
	}
	body := make([]byte, length)
	if _, err = io.ReadFull(b, body); err != nil {
		return nil, err
	}
	req.body = body
	return req, nil
}
