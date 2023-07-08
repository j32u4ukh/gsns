package account

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"strconv"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
)

// TODO: HTTP 請求處理過程中若失敗，要返回錯誤訊息給客戶端，而非印出日誌或直接返回
func (m *AccountMgr) HttpHandler(router *ans.Router) {

	// TODO: 註冊：帳密以及個人資訊。儲存以 SHA256 加密後的密碼，而非儲存原始密碼。
	router.POST("/register", m.register)

	// TODO: token 應該要有時效
	router.POST("/login", m.login)

	// TODO:
	router.POST("/logout", m.logout)

	// 取得用戶資訊
	router.GET("/user_info", m.getUserInfo)
	router.POST("/user_info", m.setUserInfo)
}

func (m *AccountMgr) register(c *ghttp.Context) {
	ap := &AccountProtocol{}
	c.ReadJson(ap)
	m.logger.Info("AccountProtocol: %+v", ap)

	// 帳號名稱(Account) 和 密碼原文(Password) 為必須，個人資訊(Info) 可以不填
	// TODO: 在前端就加密
	if ap.Account == "" || ap.Password == "" {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": fmt.Sprintf("缺少參數, account: %+v", ap),
		})
		m.httpAnswer.Send(c)
		return
	}

	// 將當前 Http 的工作結束
	m.httpAnswer.Finish(c)

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = int32(define.CommissionCommand)
	agreement.Service = int32(define.Register)
	agreement.Cid = c.GetId()

	account := &pbgo.Account{
		Account:  ap.Account,
		Password: ap.Password,
		Info:     ap.Info,
	}

	// 帳號名稱
	account.Account = ap.Account

	// 密碼原文()
	account.Password = ap.Password

	account.Info = ap.Info
	m.logger.Info("Register account: %+v", account)
	agreement.Accounts = append(agreement.Accounts, account)

	// ==================================================
	// 準備將請求轉送給 Account server
	// ==================================================
	// 形成 "建立使用者" 的請求
	td := base.NewTransData()

	// 寫入 agreement
	bs, _ := agreement.Marshal()
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
	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = int32(define.CommissionCommand)
	agreement.Service = int32(define.Login)
	agreement.Cid = c.GetId()

	td := base.NewTransData()

	ap := &AccountProtocol{}
	c.ReadJson(ap)
	m.logger.Info("AccountProtocol: %+v", ap)

	if ap.Account == "" || ap.Password == "" {
		var ret string
		if ap.Account == "" && ap.Password == "" {
			ret = "Account & password are empty."
		} else if ap.Account == "" {
			ret = "Account is empty."
		} else {
			ret = "Password is empty."
		}
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": ret,
		})
		m.httpAnswer.Send(c)
		return
	}

	account := &pbgo.Account{
		Account:  ap.Account,
		Password: ap.Password,
	}

	agreement.Accounts = append(agreement.Accounts, account)

	// 寫入 agreement
	bs, _ := agreement.Marshal()
	td.AddByteArray(bs)
	data := td.FormData()

	// 將登入數據傳到 Account 伺服器
	err := gos.SendToServer(define.AccountServer, &data, td.GetLength())

	if err != nil {
		m.logger.Error("Failed to send to server %d: %v\nError: %+v", define.DbaServer, data, err)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": "Failed to send to server.",
		})
		m.httpAnswer.Send(c)
		return
	} else {
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
	// TODO: 通知 Account server 將用戶登出。
}

func (m *AccountMgr) getUserInfo(c *ghttp.Context) {
	defer m.httpAnswer.Send(c)
	var sToken string
	var ok bool

	if sToken, ok = c.Params["token"]; !ok {
		return
	}

	token, err := strconv.ParseUint(sToken, 10, 64)

	if err != nil {
		return
	}

	user, ok := m.users.GetByKey2(token)

	if ok {
		c.Json(ghttp.StatusOK, ghttp.H{
			"name": user.Name,
			"info": user.Info,
		})
	} else {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"msg": fmt.Sprintf("Not found token %d", token),
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

	account := &pbgo.Account{
		Index:   user.Index,
		Account: user.Name,
		Info:    user.Info,
	}
	m.logger.Info("account: %+v", account)

	agreement.Accounts = append(agreement.Accounts, account)

	// 形成 "更新用戶資訊" 的請求
	td := base.NewTransData()
	// 寫入 agreement
	bs, _ := agreement.Marshal()
	td.AddByteArray(bs)
	data := td.FormData()

	// 將新用戶資訊數據傳到 Account 伺服器
	err := gos.SendTransDataToServer(define.AccountServer, td)

	if err != nil {
		fmt.Printf("(s *MainServer) CommissionHandler | Failed to send to server %d: %v\nError: %+v\n", define.DbaServer, data, err)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": "Failed to send to server",
		})
		m.httpAnswer.Send(c)
		return
	}

	m.httpAnswer.Finish(c)
}
