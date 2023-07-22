package main

import (
	"fmt"
	"net"
)

var cnt int = 1

type User struct {
	Name   string
	Addr   string
	Ch     chan string
	conn   net.Conn
	server *Server
}

// Create a user API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := User{
		Name:   fmt.Sprintf("user%d", cnt),
		Addr:   userAddr,
		Ch:     make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	cnt++
	return &user
	//Start a goroutine that listens for messages on the current User Channel

}

// Monitor the User's channel for messages and send the messages to the peer client
func (this *User) ListenMessage() {
	for {
		meg := <-this.Ch
		this.conn.Write([]byte(meg + "\n"))
	}

}

// user online
func (this *User) Online() {
	//Add User to the list
	this.server.maplock.Lock()
	this.server.OnlineMay[this.Name] = this
	this.server.maplock.Unlock()

	//Broadcast User Online Message
	this.server.BroadCast(this, "上线")
}

// user offline
func (this *User) Offline() {
	//Remove User from the list
	this.server.maplock.Lock()
	delete(this.server.OnlineMay, this.Name)
	this.server.maplock.Unlock()

	//Broadcast User Offline Message
	this.server.BroadCast(this, "下线")

}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string) {
	if msg == "who" {
		//Query all online users
		this.server.maplock.Lock()
		for _, user := range this.server.OnlineMay {
			onlineMeg := fmt.Sprintf("[%v]%v:在线\n", user.Name, user.Addr)
			this.SendMsg(onlineMeg)
		}
		this.server.maplock.Unlock()
	} else {
		this.server.BroadCast(this, msg)
	}
}
