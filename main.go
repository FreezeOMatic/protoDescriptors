package main

import (
	_ "rpcdescriptors/gen"
	"rpcdescriptors/gen/client"
	"rpcdescriptors/server"
)

func main() {
	msg := server.GetMessage()
	client.ServeMessage(msg)
}
