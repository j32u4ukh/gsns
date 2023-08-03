package gsns

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
	"github.com/j32u4ukh/gos/utils"
)

// TODO: HTTP 請求處理過程中若失敗，要返回錯誤訊息給客戶端，而非印出日誌或直接返回
// [endpoint]/social
func (s *MainServer) HttpSocialHandler(router *ans.Router) {
	// 取得其他用戶的清單
	router.GET("/other_users", s.getOtherUsers)
	router.POST("/subscribe", s.subscribe)
}

// [endpoint]/social/other_users
func (s *MainServer) getOtherUsers(c *ghttp.Context) {
	var sToken string
	var ok bool

	if sToken, ok = c.Params["token"]; !ok {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": "Not found parameter token",
		})
		s.HttpAnswer.Send(c)
		return
	}

	token, err := strconv.ParseUint(sToken, 10, 64)

	if err != nil {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": "Invalid token",
		})
		s.HttpAnswer.Send(c)
		return
	}

	user, ok := s.AMgr.GetUserByToken(token)

	if !ok {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 3,
			"msg": fmt.Sprintf("Not found user with token(%d)", token),
		})
		s.HttpAnswer.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.GetOtherUsers
	agreement.Cid = c.GetId()
	account := &pbgo.Account{
		Index: user.Index,
	}
	utils.Info("Current user id: %d", user.Index)
	agreement.Accounts = append(agreement.Accounts, account)
	bs, _ := agreement.Marshal()

	// 寫入 agreement
	td := base.NewTransData()
	td.AddByteArray(bs)
	data := td.FormData()
	logger.Info("data: %+v", data)

	// 將註冊數據傳到 Account 伺服器
	err = gos.SendToServer(define.AccountServer, &data, int32(len(data)))

	if err != nil {
		utils.Error("Failed to send request to account server, err: %+v", err)
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 4,
			"msg": "Failed to send request to account server",
		})
		s.HttpAnswer.Send(c)
		return
	}
	s.HttpAnswer.Finish(c)
}

// [endpoint]/social/subscribe
func (s *MainServer) subscribe(c *ghttp.Context) {
	ip := &SocialProtocol{}
	c.ReadJson(ip)

	if ip.Token == 0 || ip.TargetId == 0 {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": "缺少參數",
		})
		s.HttpAnswer.Send(c)
		return
	}

	user, ok := s.AMgr.GetUserByToken(ip.Token)

	if !ok {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": fmt.Sprintf("Not found token %d", ip.Token),
		})
		s.HttpAnswer.Send(c)
		return
	}

	// TODO: 應避免重複訂閱，先檢查訂閱對象的 ID 再送出請求

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.Subscribe
	agreement.Cid = c.GetId()
	agreement.Edges = append(agreement.Edges, &pbgo.Edge{
		UserId: user.Index,
		Target: ip.TargetId,
	})

	// 寫入 agreement
	bs, _ := agreement.Marshal()
	td := base.NewTransData()
	td.AddByteArray(bs)
	data := td.FormData()
	logger.Info("data: %+v", data)

	// 將註冊數據傳到 Account 伺服器
	err := gos.SendToServer(define.AccountServer, &data, int32(len(data)))

	if err != nil {
		utils.Error("Failed to send request to account server, err: %+v", err)
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 3,
			"msg": "Failed to send request to account server",
		})
		s.HttpAnswer.Send(c)
		return
	}
	s.HttpAnswer.Finish(c)
}
