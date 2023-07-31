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
	// 讀取貼文
	GetPost
	// 讀取貼文
	GetMyPosts
	// 編輯貼文
	ModifyPost
	// 自我介紹
	Introduction
	// 取得其他用戶的列表
	GetOtherUsers
	// 訂閱其他用戶
	Subscribe
)
