package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// ===== 轉 byte 陣列 =====
// 數字 轉 byte 陣列
func NumberToBytes[T int8 | int16 | int32 | int64 | uint16 | uint32 | uint64 | float32 | float64](v T, order binary.ByteOrder) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, order, v)
	return bytesBuffer.Bytes()
}

// ===== byte 陣列轉回原始數值 =====

func BytesToNumber[T int8 | int16 | int32 | int64 | uint16 | uint32 | uint64 | float32 | float64](b []byte, order binary.ByteOrder) T {
	var result T
	buffer := bytes.NewBuffer(b)
	binary.Read(buffer, order, &result)
	return result
}

func main() {
	Snowflake := make([]byte, 8)
	now := time.Now()
	timestamp := now.UnixMilli()
	fmt.Printf("now: %+v, timestamp: %d\n", now, timestamp)
	str := fmt.Sprintf("%d", timestamp)
	fmt.Printf("len(timestamp) = %d\n", len(str))
	timeBs := NumberToBytes(timestamp, binary.LittleEndian)
	fmt.Printf("timeBs(%d): %+v\n", len(timeBs), timeBs)
	copy(Snowflake[2:], timeBs[:6])
	fmt.Printf("Snowflake(%d): %+v\n", len(Snowflake), Snowflake)
	number := BytesToNumber[uint64](timeBs, binary.LittleEndian)
	fmt.Printf("number: %d\n", number)
}
