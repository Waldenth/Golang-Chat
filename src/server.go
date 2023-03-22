package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建一个server
// 以一个大写字母开头,使用这种形式的标识符的对象就可以被外部包的代码所使用
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

// 链接处理业务
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("建立连接成功")
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
