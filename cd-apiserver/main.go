package main

import (
	"pulumigcp/server"
)

func main() {
	s := server.NewServer()
	s.Run(":1323")
}
