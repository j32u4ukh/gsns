package dba

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/database"
)

type DbaServer struct {
	db     *database.Database
	DbName string
	tables map[int]*gosql.Table
	// key: server id, value: conn id
	serverIdDict map[int32]int32
}

func NewDbaServer() *DbaServer {
	s := &DbaServer{
		db:           nil,
		DbName:       "",
		tables:       make(map[int]*gosql.Table),
		serverIdDict: make(map[int32]int32),
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
		s.handleSystem(work, agreement)
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

func (s *DbaServer) handleSystem(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	// 回應心跳包
	case define.Heartbeat:
		bs, _ := agreement.Marshal()
		work.Body.AddByteArray(bs)
		work.SendTransData()
	case define.Introduction:
		if agreement.Cipher != define.CIPHER {
			logger.Error("Cipher: %s, Identity: %d", agreement.Cipher, agreement.Identity)
			gos.Disconnect(define.DbaPort, work.Index)
		} else {
			s.serverIdDict[agreement.Identity] = work.Index
			logger.Info("Hello %s from %d", define.ServerName(agreement.Identity), work.Index)
		}
		work.Finish()
	default:
		fmt.Printf("Unsupport service: %d\n", agreement.Service)
		work.Finish()
	}
}

func (s *DbaServer) handleNormalCommand(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.GetUserData:
		logger.Debug("GetUserData")
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
	case define.GetPost:
		// 只有 Account: 取得這些帳號的所有貼文
		// 有 PostMessage 列表: 取得這些 post_id 的貼文
		if len(agreement.Accounts) > 0 {
			work.Finish()
			userIds := []any{}
			for _, account := range agreement.Accounts {
				userIds = append(userIds, account.Index)
			}
			selector := s.tables[TidPostMessage].GetSelector()
			defer s.tables[TidPostMessage].PutSelector(selector)
			selector.SetCondition(gosql.WS().In("user_id", userIds...))
			pms, err := selector.Query(func() any { return &pbgo.PostMessage{} })
			if err != nil {
				agreement.ReturnCode = 1
				agreement.Msg = fmt.Sprintf("Failed to query posts from %+v.", userIds)
				logger.Error("Failed to query posts, error: %v", err)
			} else {
				agreement.ReturnCode = 0
				for _, pm := range pms {
					agreement.PostMessages = append(agreement.PostMessages, pm.(*pbgo.PostMessage))
				}
			}

			// 將用戶貼文送往 PostMessage server
			td := base.NewTransData()
			bs, _ := agreement.Marshal()
			td.AddByteArray(bs)
			data := td.FormData()
			err = gos.SendToClient(define.DbaPort, s.serverIdDict[define.PostMessageServer], &data, int32(len(data)))
			if err != nil {
				logger.Error("Failed to send to PostMessage server, err: %+v", err)
			}
		} else if len(agreement.PostMessages) > 0 {
			postIds := []any{}
			for _, pm := range agreement.PostMessages {
				postIds = append(postIds, pm.Id)
			}
			selector := s.tables[TidPostMessage].GetSelector()
			defer s.tables[TidPostMessage].PutSelector(selector)
			selector.SetCondition(gosql.WS().In("id", postIds))
			pms, err := selector.Query(func() any { return &pbgo.PostMessage{} })
			if err != nil {
				agreement.ReturnCode = 2
				agreement.Msg = "Failed to query posts."
			} else {
				agreement.ReturnCode = 0
				for _, pm := range pms {
					agreement.PostMessages = append(agreement.PostMessages, pm.(*pbgo.PostMessage))
				}
			}
			bs, _ := agreement.Marshal()
			work.Body.AddByteArray(bs)
			work.SendTransData()
		} else {
			logger.Error("Undefine which posts to query.")
			work.Finish()
		}
	}
}

func (s *DbaServer) handleCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.Register:
		// work.Body.Clear()
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
		// work.Body.Clear()
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

	case define.AddPost:
		// work.Body.Clear()
		defer func() {
			bs, _ := agreement.Marshal()
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}()

		post := agreement.PostMessages[0]
		// 根據雪花算法，給出 post id
		post.Id = GetSnowflake(0, 0)
		logger.Info("post: %+v", post)
		inserter := s.tables[TidPostMessage].GetInserter()
		defer s.tables[TidPostMessage].PutInserter(inserter)
		err := inserter.Insert(post)
		if err != nil {
			fmt.Printf("Insert err: %+v", err)
			agreement.ReturnCode = 1
			agreement.Msg = "Failed to insert account."
			agreement.PostMessages = agreement.PostMessages[:0]
			return
		}

		var result *database.SqlResult
		result, err = inserter.Exec()

		if err != nil {
			fmt.Printf("Insert Exec err: %+v\n", err)
			agreement.ReturnCode = 2
			agreement.Msg = "Failed to execute insert statement."
			agreement.PostMessages = agreement.PostMessages[:0]
			return
		}

		logger.Info("result: %s, post: %+v", result, agreement.PostMessages[0])
		agreement.ReturnCode = 0

	case define.GetPost:
		defer func() {
			bs, _ := agreement.Marshal()
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}()
		pm := agreement.PostMessages[0]
		selector := s.tables[TidPostMessage].GetSelector()
		defer s.tables[TidPostMessage].PutSelector(selector)
		selector.SetCondition(gosql.WS().Eq("id", pm.Id))
		pms, err := selector.Query(func() any { return &pbgo.PostMessage{} })
		if err != nil {
			agreement.ReturnCode = 2
			agreement.Msg = fmt.Sprintf("Failed to query post(%d).", pm.Id)
		} else if len(pms) == 0 {
			agreement.ReturnCode = 3
			agreement.Msg = fmt.Sprintf("Not found post with id(%d).", pm.Id)
		} else {
			agreement.ReturnCode = 0
			agreement.PostMessages[0] = pms[0].(*pbgo.PostMessage)
		}

	case define.ModifyPost:
		defer func() {
			bs, _ := agreement.Marshal()
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}()
		pm := agreement.PostMessages[0]
		updater := s.tables[TidPostMessage].GetUpdater()
		defer s.tables[TidPostMessage].PutUpdater(updater)
		updater.UpdateAny(pm)
		updater.SetCondition(gosql.WS().Eq("id", pm.Id))
		result, err := updater.Exec()
		if err != nil {
			agreement.ReturnCode = 1
			agreement.Msg = fmt.Sprintf("Failed to modify post(%d).", pm.Id)
		} else {
			logger.Info("Modify result: %+v", result)
			agreement.ReturnCode = 0
		}

	case define.GetOtherUsers:
		defer func() {
			bs, _ := agreement.Marshal()
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}()

		requester := agreement.Accounts[0].Index
		selector := s.tables[TidAccount].GetSelector()
		defer s.tables[TidAccount].PutSelector(selector)
		// selector.SetSelectItem(stmt.NewSelectItem("id"))
		logger.Info("requester: %d", requester)
		selector.SetCondition(gosql.WS().Ne("index", requester))
		results, err := selector.Query(func() any { return &pbgo.Account{} })
		if err != nil {
			agreement.ReturnCode = 1
			agreement.Msg = "Failed to select other users' list."
			logger.Error("GetOtherUsers err: %+v", err)
			return
		}
		var account *pbgo.Account
		agreement.Accounts = agreement.Accounts[:0]
		for _, result := range results {
			account = result.(*pbgo.Account)
			account.Password = ""
			account.CreateTime = nil
			account.UpdateTime = nil
			logger.Debug("account: %+v", account)
			agreement.Accounts = append(agreement.Accounts, account)
		}

		agreement.ReturnCode = 0

	case define.Subscribe:
		defer func() {
			bs, _ := agreement.Marshal()
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}()

		edge := agreement.Edges[0]
		inserter := s.tables[TidEdge].GetInserter()
		defer s.tables[TidEdge].PutInserter(inserter)
		err := inserter.Insert(edge)

		if err != nil {
			agreement.ReturnCode = 1
			agreement.Msg = "Failed to insert edge."
			agreement.Edges = agreement.Edges[:0]
			logger.Error("Subscribe err: %+v", err)
			return
		}

		var result *database.SqlResult
		result, err = inserter.Exec()

		if err != nil {
			fmt.Printf("Insert Exec err: %+v\n", err)
			agreement.ReturnCode = 2
			agreement.Msg = "Failed to execute insert statement."
			agreement.Edges = agreement.Edges[:0]
			return
		}

		logger.Info("result: %s, edge: %+v", result, edge)
		agreement.ReturnCode = 0

	default:
		fmt.Printf("Unsupport commission: %d\n", agreement.Service)
		work.Finish()
	}
}
