package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播channel
	Message chan string
}

// 创建一个server
// 以一个大写字母开头,使用这种形式的标识符的对象就可以被外部包的代码所使用
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听Message广播channel的goroutine,一旦有消息发送给全部在线User
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		// 将msg发送给全部在线user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg + "\n"

	this.Message <- sendMsg
}

// 处理新连接业务
func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("Connection established successfully.")
	// 新的连接,创建新的用户对象
	user := NewUser(conn, this)

	user.Online()

	// 监听用户是否活跃的channel
	isAlive := make(chan bool)

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				_, ok := user.server.OnlineMap[user.Name]
				if ok {
					user.Offline()
				}
				return
			}
			if err != nil {
				fmt.Println("Conn Read err:", err)
				return
			}

			// 提取客户端发送的信息
			msg := string(buf[:n-1])

			// 用户针对msg处理,如何发送
			user.DoMessage(msg)

			// 用户任意消息代表用户活跃
			isAlive <- true
		}
	}()

	// 阻塞当前handler 超时检测1min
	for {
		select { // 扫描全部case,重置定时器
		case <-isAlive:
			//当前用户活跃,循环下一次select,重置定时器
		case <-time.After(time.Second * 60): //扫描case,(不是进case)就会重置定时器time.After(time.Second * 10),重新倒数
			//客户端已经超时
			//强制关闭客户端
			//销毁资源
			user.SendMsg("Your connection is closed : Long time no response.\n")
			close(user.C)

			conn.Close()
			// 退出当前处理goroutine
			return
		}
	}

}

// 启动服务器接口
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen.err:", err)
	}
	//close listen socket
	defer listener.Close()

	// 启动监听message 的goroutine
	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}

}
