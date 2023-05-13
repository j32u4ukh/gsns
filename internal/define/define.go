package define

// Server id
const (
	DbaServer     = 0
	AccountServer = 1
)

// TODO: Connection config
const (
	MainPort    = 1023
	DbaPort     = 1022
	AccountPort = 1021
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
	Heartbeat uint16 = iota
	GetUserData
	Register
	Login
)
