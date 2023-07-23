package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	mod        int
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server ip address(default: 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "set server port(default: 8888)")
	flag.Parse()
}

func NewClient(ServerIp string, ServerPort int) *Client {
	//new client
	client := &Client{
		ServerIp:   ServerIp,
		ServerPort: ServerPort,
		mod:        -1,
	}

	//connect server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, ServerPort))

	if err != nil {
		fmt.Println("net.Dial err=", err)
		return nil
	}

	client.conn = conn

	return client
}

func (this *Client) menu() bool {
	var mod int

	fmt.Println("==========menu==========")
	fmt.Println("1. public chat")
	fmt.Println("2. private chat")
	fmt.Println("3. rename")
	fmt.Println("4. list all user")
	fmt.Println("0. exit")
	fmt.Println("please choose(0-3):")
	fmt.Println("========================")

	fmt.Scanln(&mod)

	if !(mod >= 0 && mod <= 4) {
		fmt.Println("====>input error, please input again<====")
		return false
	}
	this.mod = mod
	return true
}

func (this *Client) Run() {
	for this.mod != 0 {
		for this.menu() != true {
		}

		switch this.mod {
		case 1: //public chat
			this.PublicChat()
		case 2: //private chat
			this.PrivateChat()
		case 3: //rename
			this.UpdateName()
		case 4: //list all user
			this.SelectUser()
		case 0: //exit
			fmt.Println("exit")
		}

	}

}

func (this *Client) PublicChat() {
	fmt.Println(">>>>>input message to public chat,use 'exit' to exit<<<<<<")
	var chatMsg string
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {

		if len(chatMsg) != 0 { //message not empty
			sendMsg := chatMsg + "\n"
			_, err := this.conn.Write([]byte(sendMsg))

			if err != nil {
				fmt.Println("conn.Write err=", err)
				break
			}
		}
		fmt.Println(">>>>>input message to public chat,use 'exit' to exit<<<<<<")
		fmt.Scanln(&chatMsg)
	}

}

func (this *Client) SelectUser() {
	sendMsg := "list-all\n"
	_, err := this.conn.Write([]byte(sendMsg))

	if err != nil {
		fmt.Println("conn.Write err=", err)
	}
}

func (this *Client) PrivateChat() {

	var chatMsg string
	var receiver string = ""
	fmt.Println("*******************" + this.Name + "*****************")

	fmt.Println(">>>>>input message to private chat,use 'exit' to exit<<<<<<")
	this.SelectUser()
	fmt.Println("please input receiver name:")

	for {
		fmt.Scanln(&receiver)

		if receiver == this.Name {
			fmt.Println("can't private chat with yourself\n please input again")
			continue
		} else {
			break
		}
	}

on:
	for receiver != "exit" {
		fmt.Println("please input message,user 'exit' to exit:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 { //message not empty
				sendMsg := "to|" + receiver + "|" + chatMsg + "\n"
				_, err := this.conn.Write([]byte(sendMsg))

				if err != nil {
					fmt.Println("conn.Write err=", err)
					break
				}

			}

			fmt.Println("please input message,user 'exit' to exit:")
			fmt.Scanln(&chatMsg)
		}
		break on
	}

}

func (this *Client) UpdateName() {
	fmt.Println(">>>>>input your name<<<<<<")
	fmt.Scanln(&this.Name)

	sendMsg := "rename|" + this.Name + "\n"

	_, err := this.conn.Write([]byte(sendMsg))

	if err != nil {
		fmt.Println("conn.Write err=", err)
		return
	}
}

// receive server response
func (this *Client) DealResponse() {
	//once client receive server response, output to stdout
	io.Copy(os.Stdout, this.conn)
}

func main() {

	Client := NewClient(serverIp, serverPort)

	if Client == nil {
		fmt.Println(">>>>>connect server failed<<<<<<")
		return
	}
	fmt.Println(">>>>>connect server success<<<<<<")
	go Client.DealResponse()

	Client.Run()
}
