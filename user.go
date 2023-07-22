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
	this.server.BroadCast(this, "ONELINE")
}

// user offline
func (this *User) Offline() {
	//Remove User from the list
	this.server.maplock.Lock()
	delete(this.server.OnlineMay, this.Name)
	this.server.maplock.Unlock()

	//Broadcast User Offline Message
	this.server.BroadCast(this, "OFFLINE")

}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string) {
	if msg == "who" {
		//Query all online users
		this.server.maplock.Lock()
		for _, user := range this.server.OnlineMay {
			onlineMeg := fmt.Sprintf("[%v]%v:ONLINE\n", user.Name, user.Addr)
			this.SendMsg(onlineMeg)
		}
		this.server.maplock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := msg[7:]
		_, ok := this.server.OnlineMay[newName]
		if ok {
			this.SendMsg("this name is used\n")
		} else {
			this.server.maplock.Lock()
			delete(this.server.OnlineMay, this.Name)
			this.server.OnlineMay[newName] = this
			this.server.maplock.Unlock()
			this.Name = newName
			this.SendMsg("you have update your name:" + this.Name + "\n")
		}

	} else {
		this.server.BroadCast(this, msg)
	}
}
