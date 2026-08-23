package main

import (
	"net"
)

type User struct {
	Name    string
	Address string
	Channel chan string
	Connect net.Conn
	server  *Server
}

func NewUser(con net.Conn, server *Server) *User {
	address := con.RemoteAddr().String()
	user := &User{
		Name:    address,
		Address: address,
		Channel: make(chan string),
		Connect: con,
		server:  server,
	}
	go user.ListenMessage()

	return user
}

func (u *User) Online() {
	u.server.OnlineMapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.OnlineMapLock.Unlock()
	u.server.BroadCast(u, "上线了")
}
func (u *User) Offline() {
	u.server.OnlineMapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.OnlineMapLock.Unlock()
	u.server.BroadCast(u, "下线了")
}

func (u *User) SendMsg(msg string) {
	u.server.BroadCast(u, msg)
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.Channel
		u.Connect.Write([]byte(msg + "\n"))
	}
}
