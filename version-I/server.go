package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server // server是指向对象的指针
}

// 启动server的接口
func (server *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.listen error:", err)
		return // 启动失败
	}
	// close listen socket
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept error:", err)
			continue // 继续循环监听
		}
		// do handler
		go server.Handler(conn)
	}

}

func (server *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	fmt.Println("-----连接建立成功-----")
}
