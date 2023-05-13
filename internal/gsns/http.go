package gsns

import (
	"encoding/json"
	"fmt"
	"internal/define"
	"internal/pbgo"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
	"google.golang.org/protobuf/proto"
)

func (s *MainServer) register(c *ghttp.Context) {
	// 形成 "建立使用者" 的請求
	td := base.NewTransData()
	td.AddByte(define.CommissionCommand)
	td.AddUInt16(define.Register)
	td.AddInt32(c.GetId())

	account := &pbgo.Account{}
	m := map[string]string{}
	err := json.Unmarshal(c.Body[:c.BodyLength], &m)

	if err != nil {
		logger.Error("Failed to unmarshal data: %+v\nError: %+v", c.Body[:c.BodyLength], err)
		s.HttpAnswer.Finish(c)
		return
	}

	logger.Info("json body: %+v", m)
	var ok bool
	var value string

	if value, ok = m["account"]; ok {
		// 帳號名稱
		account.Account = value
	} else {
		logger.Error("Not found param: account")
		// 將當前 Http 的工作結束
		s.HttpAnswer.Finish(c)
		return
	}

	if value, ok = m["password"]; ok {
		// 密碼原文(TODO: 在前端就加密?)
		account.Password = value
	} else {
		logger.Error("Not found param: password")

		// 將當前 Http 的工作結束
		s.HttpAnswer.Finish(c)
		return
	}

	// 將當前 Http 的工作結束
	s.HttpAnswer.Finish(c)

	logger.Info("Register account: %+v", account)

	// 寫入 pbgo.Account
	bs, _ := proto.Marshal(account)
	td.AddByteArray(bs)

	data := td.FormData()

	logger.Info("data: %+v", data)

	// 將註冊數據傳到 Account 伺服器
	err = gos.SendToServer(EAccountServer, &data, td.GetLength())

	if err != nil {
		fmt.Printf("(s *MainServer) CommissionHandler | Failed to send to server %d: %v\nError: %+v\n", EDbaServer, data, err)
		return
	}
}

func (s *MainServer) login(c *ghttp.Context) {
	td := base.NewTransData()
	td.AddByte(define.CommissionCommand)
	td.AddUInt16(define.Login)
	td.AddInt32(c.GetId())

	m := map[string]string{}
	err := json.Unmarshal(c.Body[:c.BodyLength], &m)

	if err != nil {
		logger.Error("Failed to unmarshal data: %+v\nError: %+v", c.Body[:c.BodyLength], err)
		s.HttpAnswer.Finish(c)
		return
	}

	logger.Info("json body: %+v", m)
	var ok bool
	var value string
	account := &pbgo.Account{}

	if value, ok = m["account"]; ok {
		// 帳號名稱
		account.Account = value
	} else {
		logger.Error("Not found param: account")
		// 將當前 Http 的工作結束
		s.HttpAnswer.Finish(c)
		return
	}

	if value, ok = m["password"]; ok {
		// 密碼原文(TODO: 在前端就加密?)
		account.Password = value
	} else {
		logger.Error("Not found param: password")

		// 將當前 Http 的工作結束
		s.HttpAnswer.Finish(c)
		return
	}

	// 將當前 Http 的工作結束
	s.HttpAnswer.Finish(c)

	// 寫入 pbgo.Account
	bs, _ := proto.Marshal(account)
	td.AddByteArray(bs)

	data := td.FormData()

	logger.Info("data: %+v", data)

	// 將登入數據傳到 Account 伺服器
	err = gos.SendToServer(EAccountServer, &data, td.GetLength())

	if err != nil {
		logger.Error("Failed to send to server %d: %v\nError: %+v", EDbaServer, data, err)
		return
	}
}

func (s *MainServer) logout(c *ghttp.Context) {
	c.Json(200, ghttp.H{
		"index": 2,
		"msg":   "POST | /",
	})
}
