package dba

import (
	"fmt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/database"
	"google.golang.org/protobuf/proto"
)

type DbaServer struct {
	db     *database.Database
	DbName string
	tables map[int]*gosql.Table
}

func NewDbaServer() *DbaServer {
	s := &DbaServer{
		db:     nil,
		DbName: "",
		tables: make(map[int]*gosql.Table),
	}
	return s
}

func (s *DbaServer) Handler(work *base.Work) {
	cmd := work.Body.PopByte()

	switch cmd {
	case define.SystemCommand:
		s.handleSystemCommand(work)
	case define.NormalCommand:
		s.handleNormalCommand(work)
	case define.CommissionCommand:
		s.handleCommission(work)
	default:
		fmt.Printf("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (s *DbaServer) Run() {

}

func (s *DbaServer) handleSystemCommand(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {
	// 回應心跳包
	case define.Heartbeat:
		fmt.Printf("Heart beat! Now: %+v\n", time.Now())
		work.Body.Clear()
		work.Body.AddByte(define.SystemCommand)
		work.Body.AddUInt16(define.Heartbeat)
		work.Body.AddString("OK")
		work.SendTransData()
	default:
		fmt.Printf("Unsupport service: %d\n", service)
		work.Finish()
	}
}

func (s *DbaServer) handleNormalCommand(work *base.Work) {
	service := work.Body.PopUInt16()
	switch service {
	case define.GetUserData:
		logger.Debug("GetUserData")
		work.Body.Clear()
		work.Body.AddByte(define.NormalCommand)
		work.Body.AddUInt16(define.GetUserData)
		defer work.SendTransData()

		selector := s.tables[TidAccount].GetSelector()
		defer s.tables[TidAccount].PutSelector(selector)
		results, err := selector.Query(func() any { return &pbgo.Account{} })
		if err != nil {
			logger.Error("Select err: %+v", err)
			work.Body.AddUInt16(1)
			return
		}
		accounts := &pbgo.AccountArray{}
		var account *pbgo.Account
		for _, result := range results {
			account = result.(*pbgo.Account)
			account.CreateTime = nil
			logger.Debug("account: %+v", account)
			accounts.Accounts = append(accounts.Accounts, account)
		}
		bs, err := proto.Marshal(accounts)
		if err != nil {
			logger.Error("Select err: %+v", err)
			work.Body.AddUInt16(2)
			return
		}
		work.Body.AddUInt16(0)
		work.Body.AddByteArray(bs)
	}
}

func (s *DbaServer) handleCommission(work *base.Work) {
	commission := work.Body.PopUInt16()
	var cid int32 = work.Body.PopInt32()
	logger.Info("commission: %d, cid: %d", commission, cid)

	switch commission {
	case 1023:
		work.Body.Clear()
		work.Body.AddByte(1)
		work.Body.AddUInt16(1023)
		work.Body.AddInt32(cid)
		work.Body.AddString("Commission completed.")
		work.SendTransData()

	case define.Register:
		// 建立使用者資料
		bs := work.Body.PopByteArray()
		account := &pbgo.Account{}
		err := proto.Unmarshal(bs, account)
		work.Body.Clear()
		work.Body.AddByte(define.CommissionCommand)
		work.Body.AddUInt16(define.Register)
		work.Body.AddInt32(cid)

		if err != nil {
			logger.Error("Unmarshal account err: %+v", err)
			// TODO: send error message back to client
		} else {
			inserter := s.tables[TidAccount].GetInserter()
			defer s.tables[TidAccount].PutInserter(inserter)
			err = inserter.Insert(account)
			if err != nil {
				fmt.Printf("Insert err: %+v", err)
				return
			}

			var result *database.SqlResult
			result, err = inserter.Exec()

			if err != nil {
				fmt.Printf("Insert Exec err: %+v\n", err)
				return
			}

			logger.Info("result: %s", result)
			// returnCode
			work.Body.AddUInt16(0)
			account.Index = int32(result.LastInsertId)
			bs, _ = proto.Marshal(account)
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}

	case define.Login:
		// 建立使用者資料
		bs := work.Body.PopByteArray()
		account := &pbgo.Account{}
		err := proto.Unmarshal(bs, account)
		work.Body.Clear()
		work.Body.AddByte(define.CommissionCommand)
		work.Body.AddUInt16(define.Login)
		work.Body.AddInt32(cid)

		if err != nil {
			logger.Error("Unmarshal account err: %+v", err)
			// TODO: send error message back to client
		} else {
			// TODO: 檢查該帳號是否存在；若存在，檢查密碼是否正確

			// returnCode
			work.Body.AddUInt16(0)

			// 使用權 token
			work.Body.AddUInt64(9527)

			// 將結果回傳
			work.SendTransData()
		}
	default:
		fmt.Printf("Unsupport commission: %d\n", commission)
		work.Finish()
	}
}
