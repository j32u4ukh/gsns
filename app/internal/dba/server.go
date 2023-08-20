package dba

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"internal/utils"

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

func (s *DbaServer) handleSystem(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	// 回應心跳包
	case define.Heartbeat:
		_, err := agrt.SendWork(work, agreement)
		if err != nil {
			_, _, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "work")
			logger.Error("%s, err: %+v", agreement.Msg, err)
		}
	case define.Introduction:
		if agreement.Cipher != define.CIPHER {
			_, _, agreement.Msg = define.ErrorMessage(define.Error.WrongConnectionIdentity, agreement.Cipher, agreement.Identity)
			logger.Error(agreement.Msg)
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
	default:
		logger.Debug("Unsupport service: %d", agreement.Service)
		work.Finish()
	}
}

func (s *DbaServer) handleCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.Register:
		defer s.responseCommission(work, agreement)
		account := agreement.Accounts[0]
		inserter := s.tables[TidAccount].GetInserter()
		defer s.tables[TidAccount].PutInserter(inserter)
		err := inserter.Insert(account)

		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.InvalidInsertData, "account")
			logger.Error("%s, err: %+v", agreement.Msg, err)
			return
		}

		var result *database.SqlResult
		result, err = inserter.Exec()
		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.FailedInsertDb, "account")
			logger.Error("%s, err: %+v", agreement.Msg, err)
			return
		}

		account.Index = int32(result.LastInsertId)

	case define.Login:
		defer s.responseCommission(work, agreement)
		account := agreement.Accounts[0]

		//////////////////////////////////////////////////
		// 開始讀取帳號資料
		//////////////////////////////////////////////////
		accountSelector := s.tables[TidAccount].GetSelector()
		defer s.tables[TidAccount].PutSelector(accountSelector)
		accountSelector.SetCondition(gosql.WS().
			AddAndCondtion(gosql.WS().Eq("account", account.Account)).
			AddAndCondtion(gosql.WS().Eq("password", account.Password)))
		results, err := accountSelector.Query(func() any { return &pbgo.Account{} })
		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.FailedSelectDb, "account")
			logger.Error("%s, err: %+v", agreement.Msg, err)
			return
		}
		if len(results) != 1 {
			agreement.ReturnCode = 2
			agreement.Msg = fmt.Sprintf("讀取的結果數量不正確, #account: %d", len(results))
			logger.Error(agreement.Msg)
			return
		}
		account = results[0].(*pbgo.Account)
		account.CreateUtc = utils.TimestampToUtc(account.CreateTime)
		account.CreateTime = nil
		account.UpdateUtc = utils.TimestampToUtc(account.UpdateTime)
		account.UpdateTime = nil
		agreement.Accounts[0] = account
		//////////////////////////////////////////////////
		// 完成帳號資料讀取
		//////////////////////////////////////////////////
		//////////////////////////////////////////////////
		// 開始讀取社群資料
		//////////////////////////////////////////////////
		edgeSelector := s.tables[TidEdge].GetSelector()
		defer s.tables[TidEdge].PutSelector(edgeSelector)
		edgeSelector.SetCondition(gosql.WS().Eq("user_id", account.Index))
		results, err = edgeSelector.Query(func() any { return &pbgo.Edge{} })
		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.FailedSelectDb, "account")
			logger.Error("%s, err: %+v", agreement.Msg, err)
			return
		}
		var edge *pbgo.Edge
		for _, result := range results {
			edge = result.(*pbgo.Edge)
			logger.Info("edge: %+v", edge)
			edge.CreateUtc = utils.TimestampToUtc(edge.CreateTime)
			edge.CreateTime = nil
			edge.UpdateUtc = utils.TimestampToUtc(edge.UpdateTime)
			edge.UpdateTime = nil
			agreement.Edges = append(agreement.Edges, edge)
		}
		agreement.ReturnCode = 0
		//////////////////////////////////////////////////
		// 完成社群資料讀取
		//////////////////////////////////////////////////
		//////////////////////////////////////////////////
		// 開始讀取貼文資料
		//////////////////////////////////////////////////
		pmSelector := s.tables[TidPostMessage].GetSelector()
		defer s.tables[TidPostMessage].PutSelector(pmSelector)
		pmSelector.SetCondition(gosql.WS().Eq("user_id", account.Index))
		results, err = pmSelector.Query(func() any { return &pbgo.PostMessage{} })
		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.FailedSelectDb, "post message")
			logger.Error("%s, err: %+v", agreement.Msg, err)
			return
		}
		agreement2 := agrt.GetAgreement()
		defer agrt.PutAgreement(agreement2)
		agreement2.Cmd = define.NormalCommand
		agreement2.Service = define.GetMyPosts
		var pm *pbgo.PostMessage
		for _, result := range results {
			pm = result.(*pbgo.PostMessage)
			pm.CreateUtc = utils.TimestampToUtc(pm.CreateTime)
			pm.CreateTime = nil
			pm.UpdateUtc = utils.TimestampToUtc(pm.UpdateTime)
			pm.UpdateTime = nil
			agreement2.PostMessages = append(agreement2.PostMessages, pm)
		}
		_, err = agrt.SendToClient(define.DbaPort, s.serverIdDict[define.PostMessageServer], agreement2)
		if err != nil {
			_, agreement2.ReturnCode, agreement2.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "to PostMessage server")
			logger.Error("%s, err: %+v", agreement2.Msg, err)
		} else {
			logger.Info("Send to PostMessage server define.Login response: %+v", agreement2)
		}

	case define.SetUserData:
		defer s.responseCommission(work, agreement)
		account := agreement.Accounts[0]
		updater := s.tables[TidAccount].GetUpdater()
		defer s.tables[TidAccount].PutUpdater(updater)
		updater.UpdateAny(account)
		updater.SetCondition(gosql.WS().Eq("index", account.Index))
		_, err := updater.Exec()
		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.FailedUpdateDb, fmt.Sprintf("account: %+v", account))
			logger.Error("%s, err: %+v", agreement.Msg, err)
			agreement.Accounts = agreement.Accounts[:0]
			return
		}

		agreement.ReturnCode = 0

	case define.AddPost:
		defer s.responseCommission(work, agreement)

		post := agreement.PostMessages[0]

		// 根據雪花算法，給出 post id
		post.Id = GetSnowflake(0, 0)
		logger.Info("post: %+v", post)

		inserter := s.tables[TidPostMessage].GetInserter()
		defer s.tables[TidPostMessage].PutInserter(inserter)
		err := inserter.Insert(post)

		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.InvalidInsertData, "account")
			logger.Error("%s, err: %+v", agreement.Msg, err)
			agreement.PostMessages = agreement.PostMessages[:0]
			return
		}

		var result *database.SqlResult
		result, err = inserter.Exec()
		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.FailedInsertDb, "account")
			logger.Error("%s, err: %+v", agreement.Msg, err)
			agreement.PostMessages = agreement.PostMessages[:0]
			return
		}

		logger.Info("result: %s, post: %+v", result, agreement.PostMessages[0])
		agreement.ReturnCode = define.Error.None

	case define.GetPost:
		defer s.responseCommission(work, agreement)

		var pm *pbgo.PostMessage
		pm = agreement.PostMessages[0]
		selector := s.tables[TidPostMessage].GetSelector()
		defer s.tables[TidPostMessage].PutSelector(selector)
		// 貼文 ID 或親貼文 ID 與 pm.Id 相符都讀取
		selector.SetCondition(gosql.WS().
			AddOrCondtion(gosql.WS().Eq("id", pm.Id)).
			AddOrCondtion(gosql.WS().Eq("parent_id", pm.Id)))
		results, err := selector.Query(func() any { return &pbgo.PostMessage{} })
		agreement.PostMessages = agreement.PostMessages[:0]
		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.FailedSelectDb, fmt.Sprintf("post(%d).", pm.Id))
			logger.Error("%s, err: %+v", agreement.Msg, err)
		} else if len(results) == 0 {
			agreement.ReturnCode = 3
			agreement.Msg = fmt.Sprintf("Not found post with id(%d).", pm.Id)
		} else {
			agreement.ReturnCode = 0
			for _, result := range results {
				pm = result.(*pbgo.PostMessage)
				pm.CreateUtc = utils.TimestampToUtc(pm.CreateTime)
				pm.CreateTime = nil
				pm.UpdateUtc = utils.TimestampToUtc(pm.UpdateTime)
				pm.UpdateTime = nil
				agreement.PostMessages = append(agreement.PostMessages, pm)
			}
		}

	case define.ModifyPost:
		defer s.responseCommission(work, agreement)

		pm := agreement.PostMessages[0]
		updater := s.tables[TidPostMessage].GetUpdater()
		defer s.tables[TidPostMessage].PutUpdater(updater)
		updater.UpdateAny(pm)
		updater.SetCondition(gosql.WS().Eq("id", pm.Id))
		result, err := updater.Exec()
		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.FailedUpdateDb, fmt.Sprintf("post(%d)", pm.Id))
			logger.Error("%s, err: %+v", agreement.Msg, err)
		} else {
			logger.Info("Modify result: %+v", result)
			agreement.ReturnCode = 0
		}

	case define.GetOtherUsers:
		defer s.responseCommission(work, agreement)

		requester := agreement.Accounts[0].Index
		selector := s.tables[TidAccount].GetSelector()
		defer s.tables[TidAccount].PutSelector(selector)
		// TODO: selector.SetSelectItem(stmt.NewSelectItem("id"))
		logger.Info("requester: %d", requester)
		selector.SetCondition(gosql.WS().Ne("index", requester))
		results, err := selector.Query(func() any { return &pbgo.Account{} })
		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.FailedSelectDb, "other users' list")
			logger.Error("%s, err: %+v", agreement.Msg, err)
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
		defer s.responseCommission(work, agreement)

		edge := agreement.Edges[0]
		inserter := s.tables[TidEdge].GetInserter()
		defer s.tables[TidEdge].PutInserter(inserter)
		err := inserter.Insert(edge)

		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.InvalidInsertData, "edge")
			logger.Error("%s, err: %+v", agreement.Msg, err)
			agreement.Edges = agreement.Edges[:0]
			return
		}

		var result *database.SqlResult
		result, err = inserter.Exec()
		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.FailedInsertDb, "edge")
			logger.Error("%s, err: %+v", agreement.Msg, err)
			agreement.Edges = agreement.Edges[:0]
			return
		}

		logger.Info("result: %s, edge: %+v", result, edge)
		agreement.ReturnCode = define.Error.None

	case define.GetSubscribedPosts:
		defer s.responseCommission(work, agreement)

		userIds := []any{}
		for _, account := range agreement.Accounts {
			userIds = append(userIds, account.Index)
		}
		selector := s.tables[TidPostMessage].GetSelector()
		defer s.tables[TidPostMessage].PutSelector(selector)

		startTime := utils.TimeToString(utils.UtcToTime(agreement.StartUtc))
		stopTime := utils.TimeToString(utils.UtcToTime(agreement.StopUtc))
		logger.Info("startTime: %s", startTime)
		logger.Info("stopTime: %s", stopTime)
		selector.SetCondition(gosql.WS().
			AddAndCondtion(gosql.WS().In("user_id", userIds...)).
			AddAndCondtion(gosql.WS().Ge("update_time", startTime)).
			AddAndCondtion(gosql.WS().Le("update_time", stopTime)))
		results, err := selector.Query(func() any { return &pbgo.PostMessage{} })
		if err != nil {
			_, agreement.ReturnCode, agreement.Msg = define.ErrorMessage(define.Error.FailedSelectDb, fmt.Sprintf("posts from %+v.", userIds))
			logger.Error("%s, err: %+v", agreement.Msg, err)
		} else {
			agreement.ReturnCode = 0
			var pm *pbgo.PostMessage
			for _, result := range results {
				pm = result.(*pbgo.PostMessage)
				logger.Info("pm: %+v", pm)
				pm.CreateUtc = utils.TimestampToUtc(pm.CreateTime)
				pm.CreateTime = nil
				pm.UpdateUtc = utils.TimestampToUtc(pm.UpdateTime)
				pm.UpdateTime = nil
				agreement.PostMessages = append(agreement.PostMessages, pm)
			}
		}

	default:
		logger.Error("Unsupport commission: %d", agreement.Service)
		work.Finish()
	}
}

func (s *DbaServer) responseCommission(work *base.Work, agreement *agrt.Agreement) {
	_, err := agrt.SendWork(work, agreement)
	if err != nil {
		_, _, agreement.Msg = define.ErrorMessage(define.Error.CannotSendMessage, "work")
		logger.Error("%s, err: %+v", agreement.Msg, err)
	} else {
		logger.Info("Send %s response(%d): %+v", define.ServiceName(agreement.Service), agreement.ReturnCode, agreement)
	}
}
