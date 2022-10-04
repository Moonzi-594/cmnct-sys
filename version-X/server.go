package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
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
	fmt.Printf("[server]服务器启动在%s的端口:%d\n", server.Ip, server.Port)
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
	user := NewUser(conn, server)

	user.LogIn() // 上线

	// 监听用户是否活跃
	isAlive := make(chan bool)

	// 接收客户端的消息
	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.LogOut() // 下线
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read error:", err)
				return
			}

			// 提取用户发送的消息，并去除"\n"
			msg := string(buf[:n-1])

			user.DoMessage(msg) // 处理消息
			isAlive <- true
		}
	}()

	// 阻塞
	for {
		select {
		case <-isAlive:
			// 当前用户活跃，重置定时器
			// 激活select，更新下方定时器，故必须写在上面
		case <-time.After(time.Second * 30):
			// 30秒无反应，强制关闭当前user
			user.sendMsg("由于长时间未活动，您已被Nexon踢出冒险岛！\n")
			// 销毁资源
			close(user.C)
			conn.Close()
			// 退出当前Handler
			runtime.Goexit() // return
		}
	}

}
