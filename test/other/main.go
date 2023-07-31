package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	nums := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	bs, _ := json.Marshal(nums)
	fmt.Printf("nums: %s\n", string(bs))
}
