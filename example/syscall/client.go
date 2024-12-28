package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

func main() {
	// 创建一个TCP socket
	sockfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("socket creation failed:", err)
		os.Exit(1)
	}
	defer syscall.Close(sockfd)
	// 得到要访问的网页的地址和端口,也就是一个TCPAddr
	serverAddr, err := net.ResolveTCPAddr("tcp", "bing.com:80")
	if err != nil {
		fmt.Println("address resolution failed:", err)
		syscall.Close(sockfd)
		os.Exit(1)
	}
	// 使用syscall.Connect和创建好的Socket连接这个地址
	err = syscall.Connect(sockfd, &syscall.SockaddrInet4{
		Port: serverAddr.Port,
		Addr: [4]byte{serverAddr.IP[0], serverAddr.IP[1], serverAddr.IP[2], serverAddr.IP[3]},
	})
	if err != nil {
		fmt.Println("connection failed:", err)
		syscall.Close(sockfd)
		os.Exit(1)
	}
	// 发送一个请求
	request := "GET / HTTP/1.1\r\nHost: bing.com\r\n\r\n"
	_, err = syscall.Write(sockfd, []byte(request))
	if err != nil {
		fmt.Println("write failed:", err)
		syscall.Close(sockfd)
		os.Exit(1)
	}
	// 处理返回的结果，这里并没有解析http response
	response := make([]byte, 1024)
	n, err := syscall.Read(sockfd, response)
	if err != nil {
		fmt.Println("read failed:", err)
		syscall.Close(sockfd)
		os.Exit(1)
	}
	// 输出返回的结果
	fmt.Println(string(response[:n]))
}
