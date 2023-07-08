package dba

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/database"
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
	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	err := agreement.Init(work)
	if err != nil {
		work.Finish()
		logger.Error("Invalid data from work struct.")
		return
	}
	switch agreement.Cmd {
	case define.SystemCommand:
		s.handleSystemCommand(work, agreement)
	case define.NormalCommand:
		s.handleNormalCommand(work, agreement)
	case define.CommissionCommand:
		s.handleCommission(work, agreement)
	default:
		fmt.Printf("Unsupport command: %d\n", agreement.Cmd)
		work.Finish()
	}
}

func (s *DbaServer) Run() {

}

func (s *DbaServer) handleSystemCommand(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	// 回應心跳包
	case define.Heartbeat:
		fmt.Printf("Heart beat! Now: %+v\n", time.Now())
		work.Body.Clear()
		bs, _ := agreement.Marshal()
		work.Body.AddByteArray(bs)
		work.SendTransData()
	default:
		fmt.Printf("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *DbaServer) handleNormalCommand(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.GetUserData:
		logger.Debug("GetUserData")
		work.Body.Clear()
		defer func() {
			bs, _ := agreement.Marshal()
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}()

		selector := s.tables[TidAccount].GetSelector()
		defer s.tables[TidAccount].PutSelector(selector)
		results, err := selector.Query(func() any { return &pbgo.Account{} })
		if err != nil {
			agreement.ReturnCode = 1
			agreement.Msg = "Failed to select data."
			logger.Error("Select err: %+v", err)
			return
		}

		var account *pbgo.Account
		for _, result := range results {
			account = result.(*pbgo.Account)
			account.CreateTime = nil
			logger.Debug("account: %+v", account)
			agreement.Accounts = append(agreement.Accounts, account)
		}

		agreement.ReturnCode = 0
	}
}

func (s *DbaServer) handleCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.Register:
		work.Body.Clear()
		defer func() {
			bs, _ := agreement.Marshal()
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}()

		account := agreement.Accounts[0]
		inserter := s.tables[TidAccount].GetInserter()
		defer s.tables[TidAccount].PutInserter(inserter)
		err := inserter.Insert(account)
		if err != nil {
			fmt.Printf("Insert err: %+v", err)
			agreement.ReturnCode = 1
			agreement.Msg = "Failed to insert account."
			return
		}

		var result *database.SqlResult
		result, err = inserter.Exec()

		if err != nil {
			fmt.Printf("Insert Exec err: %+v\n", err)
			agreement.ReturnCode = 2
			agreement.Msg = "Failed to execute insert statement."
			return
		}

		logger.Info("result: %s", result)
		// returnCode
		agreement.ReturnCode = 0
		account.Index = int32(result.LastInsertId)

	case define.SetUserData:
		work.Body.Clear()
		defer func() {
			bs, _ := agreement.Marshal()
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}()

		account := agreement.Accounts[0]
		updater := s.tables[TidAccount].GetUpdater()
		defer s.tables[TidAccount].PutUpdater(updater)
		updater.UpdateAny(account)
		updater.SetCondition(gosql.WS().Eq("index", account.Index))
		sr, err := updater.Exec()
		logger.Info("Update result: %+v, err: %+v", sr, err)

		selector := s.tables[TidAccount].GetSelector()
		defer s.tables[TidAccount].PutSelector(selector)
		selector.SetCondition(gosql.WS().Eq("index", account.Index))
		objs, _ := selector.Query(func() any { return &pbgo.Account{} })
		agreement.Accounts[0] = objs[0].(*pbgo.Account)
		logger.Info("Update result: %+v", agreement.Accounts[0])

		// returnCode
		agreement.ReturnCode = 0

	default:
		fmt.Printf("Unsupport commission: %d\n", agreement.Service)
		work.Finish()
	}
}
