package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User `info:"在线用户表"`

	OnlineMapLock sync.RWMutex

	Message chan string `info:"总消息,负责事件分发"`
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (s *Server) ListenMessage() {
	fmt.Println("开始监听消息....")
	for msg := range s.Message {
		s.OnlineMapLock.RLock()
		for _, cli := range s.OnlineMap {
			cli.Channel <- msg
		}
		s.OnlineMapLock.RUnlock()
	}
}

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf(
		"[%s] %s: %s",
		time.Now().Format("15:04:05"),
		user.Name,
		msg,
	)
	s.Message <- sendMsg
}

func (s *Server) Handler(connect net.Conn) {
	user := NewUser(connect, s)
	user.Online()
	isLive := make(chan bool)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := user.Connect.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err == io.EOF {
				fmt.Println("Connect Reader err:", err)
				return
			}
			msg := string(buf[:n-1])
			user.DoMessage(msg)
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
		case <-time.After(time.Second * 10):
			user.backMsg("你已被踢出下线")
			user.Connect.Close()
			close(user.Channel)
			return
		}
	}
}

func (s *Server) Start() {

	// socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("socket监听失败...\n", err)
		// panic(err)
		return
	}
	fmt.Printf("服务启动成功,运行在%s:%d \n", s.Ip, s.Port)
	// close listen
	defer func() {
		listen.Close() // 关闭
	}()

	// 监听
	go s.ListenMessage()

	// accept
	for {
		connect, err := listen.Accept()
		if err != nil {
			fmt.Println("连接失败... \n", err)
			// panic(err)
			return
		}
		go func() {
			s.Handler(connect)
		}()
	}
}
