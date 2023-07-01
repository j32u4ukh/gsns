package account

import (
	"fmt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"google.golang.org/protobuf/proto"
)

type AccountServer struct {
	Tcp          *ans.Tcp0Anser
	MainServerId int32
	// key1: user index, key2: account name
	// key1 不可變更，但 key2 可以更新
	accounts *cntr.BikeyMap[int32, string, *pbgo.Account]
}

func NewAccountServer() *AccountServer {
	s := &AccountServer{
		accounts: cntr.NewBikeyMap[int32, string, *pbgo.Account](),
	}
	return s
}

func (s *AccountServer) Handler(work *base.Work) {
	cmd := work.Body.PopByte()
	logger.Info("cmd: %d", cmd)

	switch cmd {
	case define.SystemCommand:
		s.handleSystemCommand(work)
	case define.NormalCommand:
		s.handleNormal(work)
	case define.CommissionCommand:
		s.handleCommission(work)
	default:
		logger.Warn("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (rrs *AccountServer) Run() {

}

func (s *AccountServer) handleSystemCommand(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {
	// 回應心跳包
	case define.Heartbeat:
		logger.Debug("Heart beat! Now: %+v\n", time.Now())
		work.Body.Clear()
		work.Body.AddByte(0)
		work.Body.AddUInt16(0)
		work.Body.AddString("OK")
		work.SendTransData()
	default:
		logger.Warn("Unsupport service: %d\n", service)
		work.Finish()
	}
}

func (s *AccountServer) handleNormal(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {

	default:
		logger.Warn("Unsupport normal service: %d\n", service)
		work.Finish()
	}
}

func (s *AccountServer) handleCommission(work *base.Work) {
	commission := work.Body.PopUInt16()
	cid := work.Body.PopInt32()
	logger.Info("commission: %d", commission)

	switch commission {
	case define.Register:
		// TODO: 伺服器之間的連線，第一次訊息中除了前導碼，還需要自我介紹。
		s.MainServerId = work.Index
		bs := work.Body.PopByteArray()
		logger.Info("MainServerId: %d, cid: %d, bs: %+v", s.MainServerId, cid, bs)
		work.Finish()

		// ==================================================
		// 準備將請求轉送給 DBA server
		// ==================================================
		td := base.NewTransData()
		td.AddByte(define.CommissionCommand)
		td.AddUInt16(define.Register)
		td.AddInt32(cid)

		// Account data for register
		td.AddByteArray(bs)

		data := td.FormData()
		logger.Info("data: %+v", data)

		// 將註冊數據傳到 Dba 伺服器
		err := gos.SendToServer(define.DbaServer, &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to server %d: %v\nError: %+v", define.DbaServer, data, err)
			return
		}

	// TODO: 檢查用戶的名稱與密碼是否正確
	// TODO: Login 改為 Normal command?
	case define.Login:
		bs := work.Body.PopByteArray()
		data := &pbgo.Account{}
		err := proto.Unmarshal(bs, data)

		work.Body.Clear()
		work.Body.AddByte(define.CommissionCommand)
		work.Body.AddUInt16(define.Login)
		defer work.SendTransData()

		if err != nil {
			// Return code
			work.Body.AddUInt16(1)
			work.Body.AddInt32(cid)
			logger.Error("Failed to unmarshal.")
			return
		}

		var account *pbgo.Account
		var ok bool

		if account, ok = s.accounts.GetByKey2(data.Account); !ok {
			// Return code
			work.Body.AddUInt16(2)
			work.Body.AddInt32(cid)
			logger.Error("Account %s not exists.", data.Account)
			return
		}

		if data.Password != account.Password {
			// Return code
			work.Body.AddUInt16(3)
			work.Body.AddInt32(cid)
			logger.Error("Password %s is not correct.", data.Password)
			return
		}

		// Return code
		work.Body.AddUInt16(0)
		work.Body.AddInt32(cid)
		bs, _ = proto.Marshal(account)
		work.Body.AddByteArray(bs)
		logger.Info("account: %+v", account)

	// 設置用戶資料
	case define.SetUserData:
		bs := work.Body.PopByteArray()
		newAccount := &pbgo.Account{}
		proto.Unmarshal(bs, newAccount)
		account, ok := s.accounts.GetByKey1(newAccount.Index)
		if !ok {
			return
		}
		// 填入原始密碼
		newAccount.Password = account.Password

		// 當前工作直接結束，無須回應
		work.Finish()

		// ==================================================
		// 更新緩存後，再將更新請求傳送給 DBA server
		// ==================================================
		// 形成 "更新用戶資訊" 的請求
		td := base.NewTransData()
		td.AddByte(define.CommissionCommand)
		td.AddUInt16(define.SetUserData)
		td.AddInt32(cid)

		// 寫入 pbgo.Account
		bs, _ = proto.Marshal(newAccount)
		td.AddByteArray(bs)

		data := td.FormData()

		logger.Info("data: %+v", data)
		logger.Info("newData: %+v", newAccount)

		// 將新用戶資訊數據傳到 Account 伺服器
		err := gos.SendTransDataToServer(define.DbaServer, td)

		if err != nil {
			fmt.Printf("(s *MainServer) CommissionHandler | Failed to send to server %d: %v\nError: %+v\n", define.DbaServer, data, err)
			return
		}

	default:
		fmt.Printf("Unsupport commission: %d\n", commission)
		work.Finish()
	}
}
