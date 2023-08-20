package utils

import (
	"bytes"
	"encoding/binary"
)

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
