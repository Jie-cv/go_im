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
	Connect    net.Conn
	flag       int
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
		flag:       999,
	}
	connect, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("连接失败... \n", err)
		// panic(err)
		return nil
	}
	client.Connect = connect
	return client
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.Connect) // 一旦有数据，就打印出来，永久阻塞监听
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip")
	flag.IntVar(&serverPort, "port", 8888, "server port")
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")

	if _, err := fmt.Scanln(&flag); err != nil {
		fmt.Println("请输入数字")
		return false
	}

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	}

	fmt.Println("请输入合法范围内的数字")
	return false
}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入新用户名")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.Connect.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("更新用户名失败... \n", err)
		return false
	}
	fmt.Println("更新用户名成功")
	return true
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println("请输入公聊消息, exit 退出公聊")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.Connect.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("发送公聊消息失败... \n", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("请输入公聊消息, exit 退出公聊")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUsers() {
	fmt.Println("当前在线用户列表:")
	sendMsg := "who\n"
	_, err := client.Connect.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("获取在线用户列表失败... \n", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteUser string
	var chatMsg string
	client.SelectUsers()
	fmt.Println("私聊模式")
	fmt.Println("请输入私聊用户名称,exit 退出私聊")
	fmt.Scanln(&remoteUser)
	for remoteUser != "exit" {
		fmt.Println("请输入私聊消息, exit 退出私聊")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteUser + "|" + chatMsg + "\n"
				_, err := client.Connect.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("发送私聊消息失败... \n", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("请输入私聊消息, exit 退出私聊")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println("请输入私聊用户名称,exit 退出私聊")
		fmt.Scanln(&remoteUser)
	}

}

func (client *Client) Run() {
	for client.flag != 0 {
		for !client.menu() {
		}

		// 根据不同的模式处理不同的业务。
		switch client.flag {
		case 1:
			// 公聊模式
			client.PublicChat()
		case 2:
			// 私聊模式
			client.PrivateChat()
		case 3:
			// 更新用户名
			client.UpdateName()
			break
		}
	}

}

func main() {
	// 解析命令行参数
	flag.Parse()

	go func() {
		fmt.Println("开始连接服务器...")
		client := NewClient(serverIp, serverPort)
		if client == nil {
			return
		}
		defer client.Connect.Close()
		fmt.Println("连接服务器成功")
		client.Run()

		go client.DealResponse()
		select {}
	}()
	for {
	}
}
