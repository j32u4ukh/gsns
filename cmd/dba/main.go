package main

import (
	"fmt"
	"internal/dba"
)

func main() {
	fmt.Println("Hello, dba!")
	dba.Init()
	dba.Run()
}
