package define

// Command type id
const (
	// 系統任務
	SystemCommand int32 = iota
	// 一般任務
	NormalCommand
	// 轉交型請求
	CommissionCommand
)

// Request / Response id
const (
	// 系統任務 SystemCommand
	Heartbeat int32 = iota
	// 一般任務 NormalCommand
	GetUserData
	// 轉交型請求 CommissionCommand
	Register
	Login
	// 設置用戶資訊
	SetUserData
	// 新增貼文
	AddPost
)
