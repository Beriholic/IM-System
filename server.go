package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMay map[string]*User //List of online users
	maplock   sync.RWMutex
	Message   chan string
}

// Create an interface for the Server
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMay: make(map[string]*User),
		Message:   make(chan string),
	}
}

// Broadcast message
func (this *Server) BroadCast(user *User, msg string) {
	sendMeg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMeg

}

// Listen to the message channel and send the message to the client
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.maplock.Lock()

		for _, cli := range this.OnlineMay {
			cli.Ch <- msg
		}

		this.maplock.Unlock()
	}
}

// Do handler
func (this *Server) Handler(conn net.Conn) {
	user := NewUser(conn, this)

	user.Online()

	//Recieve client message
	go func() {
		//buf is used to store the message sent by the client
		//4096 is the size of the buffer
		//revieve the message sent by the client
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)

			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//Get rid of the newline('\n') character
			msg := string(buf[:n-1])

			//Broadcast message

			user.DoMessage(msg)
		}
	}()

}

// Interface for starting the server
func (this *Server) Start() {
	//socket listen
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net.Listener err:", err)
		return
	}
	//layered call close Listener
	defer Listener.Close()

	//start a goroutine to listen for channel messages
	go this.ListenMessage()

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
