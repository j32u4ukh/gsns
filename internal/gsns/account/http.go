package account

import (
	"fmt"
	"internal/define"
	"internal/pbgo"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
	"google.golang.org/protobuf/proto"
)

func (m *AccountMgr) HttpHandler(router *ans.Router) {

	// TODO: 註冊：帳密以及個人資訊。儲存以 SHA256 加密後的密碼，而非儲存原始密碼。
	router.POST("/register", m.register)

	// TODO: token 應該要有時效
	router.POST("/login", m.login)

	// TODO:
	router.POST("/logout", m.logout)

	// 轉交工作範例
	router.POST("/delay_response", func(c *ghttp.Context) {
		m.httpAnswer.Finish(c)
		// m.CommissionHandler(1023, c.GetId())
	})
}

func (m *AccountMgr) register(c *ghttp.Context) {
	// 形成 "建立使用者" 的請求
	td := base.NewTransData()
	td.AddByte(define.CommissionCommand)
	td.AddUInt16(define.Register)
	td.AddInt32(c.GetId())

	ap := &AccountProtocol{}
	c.ReadJson(ap)
	m.logger.Info("AccountProtocol: %+v", ap)
	account := &pbgo.Account{}

	if ap.Account != "" {
		// 帳號名稱
		account.Account = ap.Account
	} else {
		m.logger.Error("Not found param: account")
		// 將當前 Http 的工作結束
		m.httpAnswer.Finish(c)
		return
	}

	if ap.Password != "" {
		// 密碼原文(TODO: 在前端就加密?)
		account.Password = ap.Password
	} else {
		m.logger.Error("Not found param: password")

		// 將當前 Http 的工作結束
		m.httpAnswer.Finish(c)
		return
	}

	// 將當前 Http 的工作結束
	m.httpAnswer.Finish(c)

	m.logger.Info("Register account: %+v", account)

	// 寫入 pbgo.Account
	bs, _ := proto.Marshal(account)
	td.AddByteArray(bs)

	data := td.FormData()

	m.logger.Info("data: %+v", data)

	// 將註冊數據傳到 Account 伺服器
	err := gos.SendToServer(define.AccountServer, &data, td.GetLength())

	if err != nil {
		fmt.Printf("(s *MainServer) CommissionHandler | Failed to send to server %d: %v\nError: %+v\n", define.DbaServer, data, err)
		return
	}
}

func (m *AccountMgr) login(c *ghttp.Context) {
	td := base.NewTransData()
	td.AddByte(define.CommissionCommand)
	td.AddUInt16(define.Login)
	td.AddInt32(c.GetId())

	ap := &AccountProtocol{}
	c.ReadJson(ap)
	m.logger.Info("AccountProtocol: %+v", ap)
	account := &pbgo.Account{}

	if ap.Account != "" {
		// 帳號名稱
		account.Account = ap.Account
	} else {
		m.logger.Error("Not found param: account")
		// 將當前 Http 的工作結束
		m.httpAnswer.Finish(c)
		return
	}

	if ap.Password != "" {
		// 密碼原文(TODO: 在前端就加密?)
		account.Password = ap.Password
	} else {
		m.logger.Error("Not found param: password")

		// 將當前 Http 的工作結束
		m.httpAnswer.Finish(c)
		return
	}

	// 將當前 Http 的工作結束
	m.httpAnswer.Finish(c)

	// 寫入 pbgo.Account
	bs, _ := proto.Marshal(account)
	td.AddByteArray(bs)

	data := td.FormData()

	m.logger.Info("data: %+v", data)

	// 將登入數據傳到 Account 伺服器
	err := gos.SendToServer(define.AccountServer, &data, td.GetLength())

	if err != nil {
		m.logger.Error("Failed to send to server %d: %v\nError: %+v", define.DbaServer, data, err)
		return
	}
}

func (m *AccountMgr) logout(c *ghttp.Context) {
	defer m.httpAnswer.Send(c)
	ap := &AccountProtocol{}
	c.ReadJson(ap)

	if ap.Token == 0 {
		m.logger.Error("Not found param: token")
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"msg": "Not found token parameter.",
		})
		return
	}

	m.users.DelByKey2(ap.Token)
	c.Json(200, ghttp.H{
		"msg": "Logout success.",
	})
}
