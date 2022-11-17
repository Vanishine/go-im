package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port uint16

	UserMap map[string]*User
	mapLock sync.RWMutex

	Message chan string
}

func NewServer(ip string, port uint16) *Server {
	return &Server{
		Ip:      ip,
		Port:    port,
		UserMap: make(map[string]*User),
		Message: make(chan string),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	defer listener.Close()

	go s.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept error:", err)
			continue
		}
		go s.ConnectionHandler(conn)
	}
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		s.mapLock.RLock()
		for _, user := range s.UserMap {
			user.C <- msg
		}
		s.mapLock.RUnlock()
	}
}

func (s *Server) ConnectionHandler(conn net.Conn) {
	defer conn.Close()

	user := NewUser(conn)

	s.mapLock.Lock()
	s.UserMap[user.Username] = user
	s.mapLock.Unlock()

	s.Broadcast(user, "online")

	select {}
}

func (s *Server) Broadcast(user *User, msg string) {
	s.Message <- fmt.Sprintf("[%s]: %s", user.Username, msg)
}
