package dba

import (
	"encoding/binary"
	"internal/utils"
	"time"
)

// Table id
const (
	TidAccount int = iota
	TidPostMessage
	TidEdge
)

func GetSnowflake(machineId byte, seqId byte) uint64 {
	snowflake := make([]byte, 8)
	snowflake[0] = seqId
	snowflake[1] = machineId
	timestamp := time.Now().UnixMilli()
	timeBytes := utils.NumberToBytes(timestamp, binary.LittleEndian)
	copy(snowflake[2:], timeBytes[:6])
	return utils.BytesToNumber[uint64](timeBytes, binary.LittleEndian)
}
