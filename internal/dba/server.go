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
	// case define.GetUserData:
	// 	logger.Debug("GetUserData")
	// 	defer func() {
	// 		bs, _ := agreement.Marshal()
	// 		work.Body.AddByteArray(bs)
	// 		work.SendTransData()
	// 	}()

	// 	selector := s.tables[TidAccount].GetSelector()
	// 	defer s.tables[TidAccount].PutSelector(selector)
	// 	results, err := selector.Query(func() any { return &pbgo.Account{} })
	// 	if err != nil {
	// 		agreement.ReturnCode = 1
	// 		agreement.Msg = "Failed to select data."
	// 		logger.Error("Select err: %+v", err)
	// 		return
	// 	}

	// 	var account *pbgo.Account
	// 	for _, result := range results {
	// 		account = result.(*pbgo.Account)
	// 		account.CreateTime = nil
	// 		logger.Debug("account: %+v", account)
	// 		agreement.Accounts = append(agreement.Accounts, account)
	// 	}

	// 	agreement.ReturnCode = 0
	// case define.GetPost:
	// 	// 只有 Account: 取得這些帳號的所有貼文
	// 	// 有 PostMessage 列表: 取得這些 post_id 的貼文
	// 	if len(agreement.Accounts) > 0 {
	// 		work.Finish()
	// 		userIds := []any{}
	// 		for _, account := range agreement.Accounts {
	// 			userIds = append(userIds, account.Index)
	// 		}
	// 		selector := s.tables[TidPostMessage].GetSelector()
	// 		defer s.tables[TidPostMessage].PutSelector(selector)
	// 		selector.SetCondition(gosql.WS().In("user_id", userIds...))
	// 		pms, err := selector.Query(func() any { return &pbgo.PostMessage{} })
	// 		if err != nil {
	// 			agreement.ReturnCode = 1
	// 			agreement.Msg = fmt.Sprintf("Failed to query posts from %+v.", userIds)
	// 			logger.Error("Failed to query posts, error: %v", err)
	// 		} else {
	// 			agreement.ReturnCode = 0
	// 			for _, pm := range pms {
	// 				agreement.PostMessages = append(agreement.PostMessages, pm.(*pbgo.PostMessage))
	// 			}
	// 		}

	// 		// 將用戶貼文送往 PostMessage server
	// 		td := base.NewTransData()
	// 		bs, _ := agreement.Marshal()
	// 		td.AddByteArray(bs)
	// 		data := td.FormData()
	// 		err = gos.SendToClient(define.DbaPort, s.serverIdDict[define.PostMessageServer], &data, int32(len(data)))
	// 		if err != nil {
	// 			logger.Error("Failed to send to PostMessage server, err: %+v", err)
	// 		}
	// 	} else if len(agreement.PostMessages) > 0 {
	// 		postIds := []any{}
	// 		for _, pm := range agreement.PostMessages {
	// 			postIds = append(postIds, pm.Id)
	// 		}
	// 		selector := s.tables[TidPostMessage].GetSelector()
	// 		defer s.tables[TidPostMessage].PutSelector(selector)
	// 		selector.SetCondition(gosql.WS().In("id", postIds))
	// 		pms, err := selector.Query(func() any { return &pbgo.PostMessage{} })
	// 		if err != nil {
	// 			agreement.ReturnCode = 2
	// 			agreement.Msg = "Failed to query posts."
	// 		} else {
	// 			agreement.ReturnCode = 0
	// 			for _, pm := range pms {
	// 				agreement.PostMessages = append(agreement.PostMessages, pm.(*pbgo.PostMessage))
	// 			}
	// 		}
	// 		bs, _ := agreement.Marshal()
	// 		work.Body.AddByteArray(bs)
	// 		work.SendTransData()
	// 	} else {
	// 		logger.Error("Undefine which posts to query.")
	// 		work.Finish()
	// 	}

	default:
		logger.Debug("Unsupport service: %d", agreement.Service)
		work.Finish()
	}
}

func (s *DbaServer) handleCommission(work *base.Work, agreement *agrt.Agreement) {
	switch agreement.Service {
	case define.Register:
		defer func() {
			bs, err := agreement.Marshal()
			if err != nil {
				logger.Error("Failed to marshal agreement, err: %+v", err)
				work.Finish()
				return
			}
			logger.Info("Send define.Register response:  %+v", agreement)
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}()

		account := agreement.Accounts[0]
		inserter := s.tables[TidAccount].GetInserter()
		defer s.tables[TidAccount].PutInserter(inserter)
		err := inserter.Insert(account)
		if err != nil {
			logger.Error("Insert err: %+v", err)
			agreement.ReturnCode = 1
			agreement.Msg = "Failed to insert account."
			return
		}

		var result *database.SqlResult
		result, err = inserter.Exec()

		if err != nil {
			logger.Error("Insert exec err: %+v", err)
			agreement.ReturnCode = 2
			agreement.Msg = "Failed to execute insert statement."
			return
		}

		agreement.ReturnCode = 0
		account.Index = int32(result.LastInsertId)

	case define.Login:
		var bs []byte
		var err error
		defer func() {
			bs, err = agreement.Marshal()
			if err != nil {
				logger.Error("Failed to marshal agreement, err: %+v", err)
				work.Finish()
				return
			}
			logger.Info("Send define.Login response: %+v", agreement)
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}()
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
			agreement.ReturnCode = 1
			agreement.Msg = "Failed to query account data."
			logger.Error(agreement.Msg)
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
			agreement.ReturnCode = 3
			agreement.Msg = "Failed to query edge data."
			logger.Error(agreement.Msg)
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
			logger.Error("Failed to query post data.")
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
		bs, err = agreement2.Marshal()
		if err != nil {
			logger.Error("Failed to marshal agreement2, err: %+v", err)
			return
		}
		td := base.NewTransData()
		td.AddByteArray(bs)
		data := td.FormData()
		err = gos.SendToClient(define.DbaPort, s.serverIdDict[define.PostMessageServer], &data, int32(len(data)))
		if err != nil {
			logger.Error("Failed to marshal agreement2, err: %+v", err)
			return
		} else {
			logger.Info("Send to PostMessage server define.Login response: %+v", agreement2)
		}

	case define.SetUserData:
		var bs []byte
		var err error
		defer func() {
			bs, err = agreement.Marshal()
			if err != nil {
				logger.Error("Failed to marshal agreement, err: %+v", err)
				work.Finish()
				return
			}
			work.Body.AddByteArray(bs)
			work.SendTransData()
			logger.Info("Send define.SetUserData response: %+v", agreement)
		}()

		account := agreement.Accounts[0]
		updater := s.tables[TidAccount].GetUpdater()
		defer s.tables[TidAccount].PutUpdater(updater)
		updater.UpdateAny(account)
		updater.SetCondition(gosql.WS().Eq("index", account.Index))
		_, err = updater.Exec()

		if err != nil {
			agreement.ReturnCode = 1
			agreement.Msg = fmt.Sprintf("Failed to update account: %+v", account)
			agreement.Accounts = agreement.Accounts[:0]
			logger.Error("%s, err: %+v", agreement.Msg, err)
			return
		}

		agreement.ReturnCode = 0

	case define.AddPost:
		var bs []byte
		var err error
		defer func() {
			bs, err = agreement.Marshal()
			if err != nil {
				logger.Error("Failed to marshal agreement, err: %+v", err)
				work.Finish()
				return
			}
			work.Body.AddByteArray(bs)
			work.SendTransData()
			logger.Info("Send define.AddPost response(%d): %+v", agreement.ReturnCode, agreement)
		}()

		post := agreement.PostMessages[0]

		// 根據雪花算法，給出 post id
		post.Id = GetSnowflake(0, 0)
		logger.Info("post: %+v", post)

		inserter := s.tables[TidPostMessage].GetInserter()
		defer s.tables[TidPostMessage].PutInserter(inserter)
		err = inserter.Insert(post)
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
			agreement.ReturnCode = 2
			agreement.Msg = "Failed to execute insert statement."
			logger.Error("%s, err: %+v", agreement.Msg, err)
			agreement.PostMessages = agreement.PostMessages[:0]
			return
		}

		logger.Info("result: %s, post: %+v", result, agreement.PostMessages[0])
		agreement.ReturnCode = 0

	case define.GetPost:
		var bs []byte
		var err error
		defer func() {
			bs, err = agreement.Marshal()
			if err != nil {
				logger.Error("Failed to marshal agreement, err: %+v", err)
				work.Finish()
				return
			}
			work.Body.AddByteArray(bs)
			work.SendTransData()
			logger.Info("Send define.GetPost response(%d): %+v", agreement.ReturnCode, agreement)
		}()
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
			agreement.ReturnCode = 2
			agreement.Msg = fmt.Sprintf("Failed to query post(%d).", pm.Id)
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
		defer func() {
			bs, err := agreement.Marshal()
			if err != nil {
				logger.Error("Failed to marshal agreement, err: %+v", err)
				work.Finish()
				return
			}
			work.Body.AddByteArray(bs)
			work.SendTransData()
			logger.Info("Send define.ModifyPost response: %+v", agreement)
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
		var bs []byte
		var err error
		defer func() {
			bs, err = agreement.Marshal()
			if err != nil {
				logger.Error("Failed to marshal agreement, err: %+v", err)
				work.Finish()
				return
			}
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
		var bs []byte
		var err error
		defer func() {
			bs, err = agreement.Marshal()
			if err != nil {
				work.Finish()
				return
			}
			work.Body.AddByteArray(bs)
			work.SendTransData()
			logger.Info("Send define.Subscribe response: %+v", agreement)
		}()

		edge := agreement.Edges[0]
		inserter := s.tables[TidEdge].GetInserter()
		defer s.tables[TidEdge].PutInserter(inserter)
		err = inserter.Insert(edge)

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

	case define.GetSubscribedPosts:
		var bs []byte
		var err error
		defer func() {
			bs, err = agreement.Marshal()
			if err != nil {
				logger.Error("Failed to marshal agreement, err: %+v", err)
				work.Finish()
				return
			}
			work.Body.AddByteArray(bs)
			work.SendTransData()
			logger.Info("Send define.GetSubscribedPosts response: %+v", agreement)
		}()

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
			agreement.ReturnCode = 1
			agreement.Msg = fmt.Sprintf("Failed to query posts from %+v.", userIds)
			logger.Error("%s, err: %v", agreement.Msg, err)
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
