package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn

	flag int // client模式
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       -1,
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

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.群聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更改用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if 0 <= flag && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>请输入合法范围的整数<<<<<")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			fmt.Println("[群聊模式]")
			break
		case 2:
			fmt.Println("[私聊模式]")
			break
		case 3:
			fmt.Println("[更改用户名]")
			break
		}
	}
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
	client.Run()
}
