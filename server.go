package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		ip, port,
	}
	return server
}

func (s *Server) Handler(connect net.Conn) {
	fmt.Println("执行callback")
}

func (s *Server) Start() {

	// socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("socket监听失败...", err)
		// panic(err)
		return
	}
	fmt.Printf("服务启动成功,运行在%s:%d", s.Ip, s.Port)
	// close listen
	defer func() {
		listen.Close() // 关闭
	}()
	// accept
	for {
		connect, err := listen.Accept()
		if err != nil {
			fmt.Println("连接失败...", err)
			// panic(err)
			return
		}
		go func() {
			s.Handler(connect)
		}()
	}
}
