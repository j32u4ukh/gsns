package define

// Server id
const (
	DbaServer int32 = 0
)

// Command type id
const (
	// 系統任務
	SystemCommand byte = iota
	// 一般任務
	NormalCommand
	// 轉交型請求
	CommissionCommand
)

// Request / Response id
const (
	Register uint16 = 0
)
