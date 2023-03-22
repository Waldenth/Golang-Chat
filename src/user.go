package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server // 该user属于哪个server处理
}

// 创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// 启动监听当前user channel的goroutine
	go user.ListenMessage()

	return user
}

// 用户上线业务
func (this *User) Online() {
	//用户上线,onlineMap更新
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	//全局广播用户上线消息
	this.server.BroadCast(this, " is online now.")
}

// 用户下线业务
func (this *User) Offline() {
	//用户下线,onlineMap更新
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	//全局广播用户下线消息
	this.server.BroadCast(this, " is offline now.")

}

//给当前user对应的客户端发消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户处理消息业务
func (this *User) DoMessage(msg string) {
	if msg == "who" { // 用户查询所有当前在线用户信息
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + " is online...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()

	} else if msg == "whoami" { // 用户查询自己的信息
		this.SendMsg("Your user name is [" + this.Name + "]\n")
		_, ok := this.server.OnlineMap[this.Name]
		if ok {
			this.SendMsg("And your state is Online.\n")
		} else {
			this.SendMsg("And your state is Offline.\n")
		}

	} else if len(msg) > 7 && msg[:7] == "rename|" { // 用户更新用户名
		// rename|<newname>
		newName := strings.Split(msg, "|")[1]

		// 判断newNme是否被其他用户使用
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("Sorry, this name has been used by other users.\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName

			this.SendMsg("Now your username has been updated to[" + this.Name + "]\n")
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}

// 监听当前User Channel，
// 一旦有消息就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
