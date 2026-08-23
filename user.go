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

// 给自己的消息
func (u *User) backMsg(msg string) {
	u.Connect.Write([]byte(msg))
}

func (u *User) DoMessage(msg string) {
	if msg == "who" {
		for _, v := range u.server.OnlineMap {
			if v.Name == u.Name {
				continue
			}
			res_msg := "[" + v.Name + " 在线....]\n"
			u.backMsg(res_msg)
		}
	} else if len(msg) > 6 && msg[:7] == "rename|" {
		new_name := msg[7:]
		u.server.OnlineMapLock.Lock()
		_, ok := u.server.OnlineMap[new_name]
		u.server.OnlineMapLock.Unlock()
		if ok {
			u.backMsg("当前用户名已存在\n")
		} else {
			u.server.OnlineMapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[new_name] = u
			u.Name = new_name
			u.server.OnlineMapLock.Unlock()
			u.backMsg("用户名修改成功,修改为:" + u.Name + "\n")
		}

	} else {
		u.server.BroadCast(u, msg)
	}
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.Channel
		u.Connect.Write([]byte(msg + "\n"))
	}
}
