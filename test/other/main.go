package main

import (
	"fmt"
	"time"
)

var TIME_LAYOUT string = "2006-01-02 15:04:05"

func main() {
	utc := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	t := time.Unix(utc, 0).UTC()
	fmt.Printf("utc: %d\n", utc)
	fmt.Printf("t: %+v\n", t)
	s := t.Format(TIME_LAYOUT)
	fmt.Printf("s: %s\n", s)
}
