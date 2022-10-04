package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	// 广播消息的channel
	Message chan string
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
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

	// 启动监听Message的goroutine
	go server.ListenMessage()

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

// 监听Message，有消息便广播
func (server *Server) ListenMessage() {
	for {
		msg := <-server.Message

		// 广播给所有在线的user
		server.mapLock.Lock()
		for _, u := range server.OnlineMap {
			u.C <- msg
		}
		server.mapLock.Unlock()
	}
}

// 广播消息
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	// fmt.Println("-----连接建立成功-----")
	user := NewUser(conn)

	// 用户上线，加入OnlineMap
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	// 广播当前用户上线消息
	server.BroadCast(user, "登录成功。")

	// 接收客户端的消息
	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)
			if n == 0 {
				server.BroadCast(user, "离线。")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read error:", err)
				return
			}

			// 提取用户发送的消息，并去除"\n"
			msg := string(buf[:n-1])
			// 广播消息
			server.BroadCast(user, msg)
		}
	}()

	// 阻塞
	select {}
}
