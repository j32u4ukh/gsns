package main

import (
	"fmt"
	"internal/define"
)

var TIME_LAYOUT string = "2006-01-02 15:04:05"

func main() {
	fmt.Printf("BadRequest: %+v\n", define.Error.BadRequest)
}
