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
	flag       int //客户端模式
}

var serverIp string
var serverPort int

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       99,
	}
	// 连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial.error:", err)
		return nil
	}
	client.conn = conn
	// 返回对象
	return client
}

// 菜单
func (client *Client) menu() bool {
	var flag int

	fmt.Println("1: Broadcast Chat")
	fmt.Println("2: Private chat")
	fmt.Println("3: Update username")
	fmt.Println("0: Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>Please input valid number.<<<<<<")
		return false
	}

}

func (client *Client) BroadcastChat() {
	fmt.Println(">>>>>>Please input content,\"exit\" for exit")
	var chatMsg string = ""
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Writer.err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>>Please input content,\"exit\" for exit")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	fmt.Println(">>>>>>Current online users:")
	client.SelectUsers()
	fmt.Println(">>>>>>Please input target user, \"exit\" for exit")
	var remoteName string
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>>Please input content,\"exit\" for exit")
		var chatMsg string
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Writer.err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>>Please input content,\"exit\" for exit")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println(">>>>>>Please input content,\"exit\" for exit")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>>Please input new username:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Writer err:", err)
		return false
	}
	return true
}

// 主业务 发消息
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		// 根据不同模式处理不同业务
		switch client.flag {
		case 1: //公聊
			fmt.Println("You choose [Broadcast Chat]...")
			client.BroadcastChat()
		case 2: //私聊
			fmt.Println("You choose [Private chat]...")
			client.PrivateChat()
		case 3: //更新用户名
			fmt.Println("You choose [Update username]...")
			client.UpdateName()
		}
	}
	_, err := client.conn.Write([]byte("exit\n"))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

// 处理server回应消息的goroutine
func (client *Client) DealResponse() {
	// 一旦client.conn有数据,直接copy到stdout,永久阻塞监听
	io.Copy(os.Stdout, client.conn)
	/*
		for{
			buf:=make([]byte, 4096)
			client.conn.Read(buf)
		}
	*/
}

// 命令行解析
// init函数先于main函数自动执行，不能被其他函数调用
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server ip(default 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "set server port(default 8888)")
}

func main() {
	// 解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>connnect to server failed.")
		return
	}
	fmt.Println(">>>>>>connnect to server success.")
	// 开启另一个goroutine处理收到的消息
	go client.DealResponse()

	// 启动客户端主业务
	client.Run()

}
