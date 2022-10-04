package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 上线
func (user *User) LogIn() {
	// 加入server的OnlineMap
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	// 广播上线消息
	user.server.BroadCast(user, "登录成功")
}

// 离线
func (user *User) LogOut() {
	// 加入server的OnlineMap
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	// 广播离线消息
	user.server.BroadCast(user, "离线")
}

// 给user对应的客户端发送消息
func (user *User) sendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

// 处理消息
func (user *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询在线用户
		user.server.mapLock.Lock()
		for _, u := range user.server.OnlineMap {
			onlineMsg := "[" + u.Addr + "]" + u.Name + ":" + "在玩冒险岛...\n"
			user.sendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 8 && msg[:8] == "rename=>" {
		// 更改用户名
		newName := msg[8:]
		// 重名判断
		_, ok := user.server.OnlineMap[newName]
		if ok {
			user.sendMsg("该用户名已被占用，请重试！\n")
		} else {
			// 更改server中的OnlineMap
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()
			// 更改client的信息
			user.Name = newName
			user.sendMsg(fmt.Sprintf("您已成功更新用户名为：%s...\n", newName))
		}
	} else if len(msg) > 6 && msg[:3] == "to|" {
		arr := strings.Split(msg, "|")
		target, content := arr[1], arr[2]

		if target == "" {
			user.sendMsg("消息格式错误！请使用\"to|枣枣|死鬼~\"的格式。\n")
			return
		}

		u, ok := user.server.OnlineMap[target]
		if !ok {
			user.sendMsg("不存在该用户！请检查后重新输入...\n")
			return
		}

		if content == "" {
			user.sendMsg("空消息！请重新发送...\n")
			return
		}
		// 目标用户接收消息
		u.sendMsg("[" + user.Name + "]" + "private-chatting:" + content + "\n")

	} else {
		user.server.BroadCast(user, msg)
	}
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	// 创建后立即启动 监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 监听当前User channel的方法，一旦有消息，立即发送给客户端
func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n"))
	}
}
