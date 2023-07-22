package client

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
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
func main() {

	Client := NewClient(serverIp, serverPort)

	if Client == nil {
		fmt.Println(">>>>>connect server failed<<<<<<")
		return
	}
	fmt.Println(">>>>>connect server success<<<<<<")

	select {}
}
