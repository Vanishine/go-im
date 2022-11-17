package main

import "net"

type User struct {
	Username, Addr string
	C              chan string
	conn           net.Conn
}

func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Username: addr,
		Addr:     addr,
		C:        make(chan string),
		conn:     conn,
	}
	go user.Listen()
	return user
}

func (u *User) Listen() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
