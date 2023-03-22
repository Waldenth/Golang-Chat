package main

import (
	"net"
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

// 用户处理消息业务
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}

// 监听当前User Channel，
// 一旦有消息就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
