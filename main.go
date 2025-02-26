package main

import (
	"rpcdescriptors/client"
	_ "rpcdescriptors/gen"
	"rpcdescriptors/server"
)

func main() {
	//msg := server.GetMessage()
	//client.ServeMessage(msg)

	msg := server.GetEnumMessage()
	client.ServeEnumMessage(msg)
}
