package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port uint16
}

func NewServer(ip string, port uint16) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept error:", err)
			continue
		}
		go s.ConnectionHandler(conn)
	}
}

func (s *Server) ConnectionHandler(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection established!")
	fmt.Println(conn)
}
