package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// Create an interface for the Server
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

// Do handler
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Connection established successfully")

}

// Interface for starting the server
func (this *Server) Start() {
	//socket listen
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net.Listener err:", err)
		return
	}
	//close listen socket
	defer Listener.Close()

	for {
		//accept
		conn, err := Listener.Accept()

		if err != nil {
			fmt.Println("Listener.Accept err:", err)
			continue
		}

		//do handler
		go this.Handler(conn)

	}

}
