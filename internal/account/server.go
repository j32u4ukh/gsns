package account

import (
	"fmt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"google.golang.org/protobuf/proto"
)

type AccountServer struct {
	Tcp          *ans.Tcp0Anser
	MainServerId int32
	accounts     map[string]*pbgo.Account
}

func NewAccountServer() *AccountServer {
	s := &AccountServer{
		accounts: make(map[string]*pbgo.Account),
	}
	return s
}

func (s *AccountServer) Handler(work *base.Work) {
	cmd := work.Body.PopByte()
	logger.Info("cmd: %d", cmd)

	switch cmd {
	case define.SystemCommand:
		s.handleSystemCommand(work)
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

func (s *AccountServer) handleCommission(work *base.Work) {
	commission := work.Body.PopUInt16()
	logger.Info("commission: %d", commission)

	switch commission {
	case 1023:
		cid := work.Body.PopInt32()
		work.Body.Clear()

		work.Body.AddByte(1)
		work.Body.AddUInt16(1023)
		work.Body.AddInt32(cid)
		work.Body.AddString("Commission completed.")
		work.SendTransData()

	case define.Register:
		// TODO: 伺服器之間的連線，第一次訊息中除了前導碼，還需要自我介紹。
		s.MainServerId = work.Index
		cid := work.Body.PopInt32()
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
	case define.Login:
		cid := work.Body.PopInt32()
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

		if account, ok = s.accounts[data.Account]; !ok {
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
		name := account.Account
		logger.Info("name: %s", name)
		work.Body.AddString("account1")
		work.Body.AddInt32(account.Index)
		logger.Info("account: %+v", account)

	default:
		fmt.Printf("Unsupport commission: %d\n", commission)
		work.Finish()
	}
}
