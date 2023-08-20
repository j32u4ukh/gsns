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
	// 取得訂閱用戶的貼文
	GetSubscribedPosts
)

func ServiceName(service int32) string {
	switch service {
	case Heartbeat:
		return "Heartbeat"
	case GetUserData:
		return "GetUserData"
	case Register:
		return "Register"
	case Login:
		return "Login"
	case SetUserData:
		return "SetUserData"
	case AddPost:
		return "AddPost"
	case GetPost:
		return "GetPost"
	case GetMyPosts:
		return "GetMyPosts"
	case ModifyPost:
		return "ModifyPost"
	case Introduction:
		return "Introduction"
	case GetOtherUsers:
		return "GetOtherUsers"
	case Subscribe:
		return "Subscribe"
	case GetSubscribedPosts:
		return "GetSubscribedPosts"
	default:
		return "Unknown"
	}
}
