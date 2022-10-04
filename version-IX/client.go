package main

import (
	"flag"
	"fmt"
	"net"
)

type Clinet struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Clinet {
	client := &Clinet{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// Dial: 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn

	return client
}

// 命令行参数
var (
	serverIp   string
	serverPort int
)

// ./client -ip x.x.x.x -port xxxx
func init() {
	// flag: 命令行解析
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "服务器IP地址")
	flag.IntVar(&serverPort, "port", 8888, "服务器端口")
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("[client]>>>>> 连接服务器失败...")
	} else {
		fmt.Println("[client]>>>>> 连接服务器成功...")
	}

	// 启动客户端业务
	select {}
}
