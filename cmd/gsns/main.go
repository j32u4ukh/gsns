package main

import (
	"fmt"
	"internal/gsns"
)

func main() {
	fmt.Println("Hello, gsns!")
	gsns.Init()
	gsns.Run()
}
