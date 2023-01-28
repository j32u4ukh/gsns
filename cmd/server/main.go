package main

import (
	"fmt"
	"internal/server"
)

func main() {
	fmt.Println("Hello, gsns!")
	server.Init()
	server.Run()
}
