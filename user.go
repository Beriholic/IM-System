package main

import (
	"fmt"
	"net"
)

var cnt int = 1

type User struct {
	Name string
	Addr string
	Ch   chan string
	conn net.Conn
}

// Create a user API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := User{
		Name: fmt.Sprintf("user%d", cnt),
		Addr: userAddr,
		Ch:   make(chan string),
		conn: conn,
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
