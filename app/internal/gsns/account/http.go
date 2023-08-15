package account

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"

	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base/ghttp"
)

// TODO: HTTP 請求處理過程中若失敗，要返回錯誤訊息給客戶端，而非印出日誌或直接返回
// [endpoint]/account
func (m *AccountMgr) HttpAccountHandler(router *ans.Router) {

	// TODO: 註冊：帳密以及個人資訊。儲存以 SHA256 加密後的密碼，而非儲存原始密碼。
	router.POST("/register", m.register)

	// TODO: token 應該要有時效
	router.POST("/login", m.login)

	// TODO:
	router.POST("/logout", m.logout)

	// 取得用戶資訊
	router.POST("/get_user_info", m.getUserInfo)
	router.POST("/set_user_info", m.setUserInfo)
}

func (m *AccountMgr) register(c *ghttp.Context) {
	ap := &AccountProtocol{}
	c.ReadJson(ap)
	m.logger.Info("AccountProtocol: %+v", ap)

	// 帳號名稱(Account) 和 密碼原文(Password) 為必須，個人資訊(Info) 可以不填
	// TODO: 在前端就加密
	if ap.Account == "" || ap.Password == "" {
		msg := fmt.Sprintf("缺少參數, account: %+v", ap)
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		m.httpAnswer.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.Register
	agreement.Cid = c.GetId()
	agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
		Account:  ap.Account,
		Password: ap.Password,
		Info:     ap.Info,
	})

	// // 寫入 agreement
	// td := base.NewTransData()
	// bs, _ := agreement.Marshal()
	// td.AddByteArray(bs)
	// data := td.FormData()

	// // 將註冊數據傳到 Account 伺服器
	// err := gos.SendToServer(define.AccountServer, &data, int32(len(data)))

	_, err := agrt.SendToServer(define.AccountServer, agreement)

	if err != nil {
		m.logger.Error("Failed to send to server %d\nError: %+v", define.AccountServer, err)
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 2,
			"msg": "Failed to send to Account server",
		})
		m.httpAnswer.Send(c)
	} else {
		m.logger.Info("Send define.Register request: %+v", agreement)
		// 將當前 Http 的工作結束
		m.httpAnswer.Finish(c)
	}
}

func (m *AccountMgr) login(c *ghttp.Context) {
	ap := &AccountProtocol{}
	c.ReadJson(ap)
	m.logger.Info("AccountProtocol: %+v", ap)

	if ap.Account == "" || ap.Password == "" {
		msg := fmt.Sprintf("缺少參數, ap: %+v", ap)
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		m.httpAnswer.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = int32(define.CommissionCommand)
	agreement.Service = int32(define.Login)
	agreement.Cid = c.GetId()
	agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
		Account:  ap.Account,
		Password: ap.Password,
	})

	// // 寫入 agreement
	// td := base.NewTransData()
	// bs, _ := agreement.Marshal()
	// td.AddByteArray(bs)
	// data := td.FormData()

	// // 將登入數據傳到 Account 伺服器
	// err := gos.SendToServer(define.AccountServer, &data, int32(len(data)))

	_, err := agrt.SendToServer(define.AccountServer, agreement)

	if err != nil {
		m.logger.Error("Failed to send to Account server, err: %+v", err)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": "Failed to send to server.",
		})
		m.httpAnswer.Send(c)
	} else {
		m.logger.Info("Send define.Login request: %+v", agreement)
		// 將當前 Http 的工作結束
		m.httpAnswer.Finish(c)
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

	user, ok := m.users.GetByKey2(ap.Token)

	if !ok {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"msg": fmt.Sprintf("Not found token %d", ap.Token),
		})
	} else {
		m.users.DelByKey2(ap.Token)
		c.Json(200, ghttp.H{
			"msg": fmt.Sprintf("User %s logout success.", user.Name),
		})
	}
}

func (m *AccountMgr) getUserInfo(c *ghttp.Context) {
	defer m.httpAnswer.Send(c)
	ap := &AccountProtocol{}
	c.ReadJson(ap)

	if ap.Token == 0 {
		m.logger.Error("Not found param: token")
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": "Not found token parameter.",
		})
		m.httpAnswer.Send(c)
		return
	}

	user, ok := m.users.GetByKey2(ap.Token)

	if ok {
		c.Json(ghttp.StatusOK, ghttp.H{
			"name": user.Name,
			"info": user.Info,
		})
	} else {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"msg": fmt.Sprintf("Not found token %d", ap.Token),
		})
	}
}

func (m *AccountMgr) setUserInfo(c *ghttp.Context) {
	ap := &AccountProtocol{}
	c.ReadJson(ap)

	if ap.Token == 0 {
		m.logger.Error("Not found param: token")
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": "Not found token parameter.",
		})
		m.httpAnswer.Send(c)
		return
	}

	user, ok := m.users.GetByKey2(ap.Token)

	if !ok {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": fmt.Sprintf("Not found token %d", ap.Token),
		})
		m.httpAnswer.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = int32(define.CommissionCommand)
	agreement.Service = int32(define.SetUserData)
	agreement.Cid = c.GetId()

	// 有數值，才更新
	if ap.Account != "" {
		user.Name = ap.Account
	}

	// 有數值，才更新
	if ap.Info != "" {
		user.Info = ap.Info
	}

	agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
		Index:   user.Index,
		Account: user.Name,
		Info:    user.Info,
	})

	// // 寫入 agreement
	// td := base.NewTransData()
	// bs, _ := agreement.Marshal()
	// td.AddByteArray(bs)
	// data := td.FormData()

	// // 將新用戶資訊數據傳到 Account 伺服器
	// err := gos.SendToServer(define.AccountServer, &data, int32(len(data)))
	_, err := agrt.SendToServer(define.AccountServer, agreement)

	if err != nil {
		m.logger.Error("Failed to send to %s, err: %+v", define.ServerName(define.AccountServer), err)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": "Failed to send to server",
		})
		m.httpAnswer.Send(c)
	} else {
		m.logger.Info("Send define.SetUserData request: %+v", agreement)
		m.httpAnswer.Finish(c)
	}
}
