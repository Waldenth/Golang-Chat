package main

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

// 启动服务器接口
func (this *Server) Start() {
	// socket lesten

	//accept

	//do handler

	//close listen socket
}
