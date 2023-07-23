package main

import (
	server "IM-System/Server"
)

func main() {
	Server := server.NewServer("127.0.0.1", 8888)
	Server.Start()
}
