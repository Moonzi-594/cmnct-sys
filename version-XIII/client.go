package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

// 处理server回应的消息，直接显示到标准输出
func (client *Client) DealResponse() {
	// 永久阻塞
	io.Copy(os.Stdout, client.conn)

	/*
		for {
			buf := make([]byte)
			client.conn.Read(buf)
			fmt.Println(buf)
		}
	*/
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

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println("请输入[public]聊天内容, exit退出...")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			_, err := client.conn.Write(([]byte(chatMsg + "\n")))
			if err != nil {
				fmt.Println("conn.Write error:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("请输入[public]聊天内容, exit退出...")
		fmt.Scanln(&chatMsg)
	}
}

// 查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()

	fmt.Println("请输入私聊的[用户名], exit退出...")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("请输入私聊[内容], exit退出...")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write error:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println("请输入私聊[内容], exit退出...")
			fmt.Scanln(&chatMsg)
		}
	}

}

func (client *Client) ReName() bool {
	fmt.Println(">>>>>请输入新用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename=>" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error:", err)
		return false
	}

	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.ReName()
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
		go client.DealResponse()
	}

	// 启动客户端业务
	client.Run()
}
