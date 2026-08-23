package main

import "net"

type User struct {
	Name    string
	Address string
	Channel chan string
	Connect net.Conn
}

func NewUser(con net.Conn) *User {
	address := con.RemoteAddr().String()
	user := &User{
		Name:    address,
		Address: address,
		Channel: make(chan string),
		Connect: con,
	}
	go user.ListenMessage()

	return user
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.Channel
		u.Connect.Write([]byte(msg + "\n"))
	}
}
