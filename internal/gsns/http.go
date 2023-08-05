package gsns

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"
	"internal/utils"
	"strconv"
	"time"

	"github.com/j32u4ukh/cntr"
	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
)

// TODO: HTTP 請求處理過程中若失敗，要返回錯誤訊息給客戶端，而非印出日誌或直接返回
// [endpoint]/social
func (s *MainServer) HttpSocialHandler(router *ans.Router) {
	// 取得其他用戶的清單
	router.GET("/other_users", s.getOtherUsers)
	router.POST("/subscribe", s.subscribe)
	router.POST("/subscribed_posts", s.getSubscribedPosts)
}

// [endpoint]/social/other_users
func (s *MainServer) getOtherUsers(c *ghttp.Context) {
	var sToken string
	var ok bool

	if sToken, ok = c.Params["token"]; !ok {
		msg := "Not found parameter token"
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	token, err := strconv.ParseUint(sToken, 10, 64)

	if err != nil {
		msg := "Invalid token"
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	user, ok := s.AMgr.GetUserByToken(token)

	if !ok {
		msg := fmt.Sprintf("Not found user with token(%d)", token)
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 3,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.GetOtherUsers
	agreement.Cid = c.GetId()
	agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
		Index: user.Index,
	})
	bs, err := agreement.Marshal()

	if err != nil {
		msg := "Failed to marshal agreement."
		logger.Error(fmt.Sprintf("%s, err: %+v", msg, err))
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 4,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	// 寫入 agreement
	td := base.NewTransData()
	td.AddByteArray(bs)
	data := td.FormData()

	// 將註冊數據傳到 Account 伺服器
	err = gos.SendToServer(define.AccountServer, &data, int32(len(data)))

	if err != nil {
		msg := "Failed to send request to account server"
		logger.Error(fmt.Sprintf("%s, err: %+v", msg, err))
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 5,
			"msg": msg,
		})
		s.Http.Send(c)
	} else {
		logger.Info("Send define.GetOtherUsers request: %+v", agreement)
		s.Http.Finish(c)
	}
}

// [endpoint]/social/subscribe
func (s *MainServer) subscribe(c *ghttp.Context) {
	ip := &SocialProtocol{}
	c.ReadJson(ip)

	if ip.Token == 0 || ip.TargetId == 0 {
		msg := "缺少參數"
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	user, ok := s.AMgr.GetUserByToken(ip.Token)

	if !ok {
		msg := fmt.Sprintf("Not found token %d", ip.Token)
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	// 避免重複訂閱，先檢查訂閱對象的 ID 再送出請求
	if edges, ok := s.AMgr.Edges[user.Index]; ok {
		if edges.Contains(ip.TargetId) {
			msg := fmt.Sprintf("User %s has subscribed user %d", user.Name, ip.TargetId)
			logger.Info(msg)
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": 0,
				"msg": msg,
			})
			s.Http.Send(c)
			return
		}
	}

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
	bs, err := agreement.Marshal()

	if err != nil {
		msg := "Failed to marshal agreement."
		logger.Error(fmt.Sprintf("%s, err: %+v", msg, err))
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 3,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	td := base.NewTransData()
	td.AddByteArray(bs)
	data := td.FormData()

	// 將註冊數據傳到 Account 伺服器
	err = gos.SendToServer(define.AccountServer, &data, int32(len(data)))

	if err != nil {
		msg := "Failed to send request to account server"
		logger.Error("%s, err: %+v", msg, err)
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 4,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	} else {
		logger.Info("Send define.Subscribe request: %+v", agreement)
		s.Http.Finish(c)
	}
}

// [endpoint]/social/subscribed_posts
func (s *MainServer) getSubscribedPosts(c *ghttp.Context) {
	ip := &SocialProtocol{}
	c.ReadJson(ip)

	if ip.Token == 0 {
		msg := "缺少參數"
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	var user *pbgo.SnsUser
	var ok bool
	user, ok = s.AMgr.GetUserByToken(ip.Token)

	if !ok {
		msg := fmt.Sprintf("Not found token %d", ip.Token)
		logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)

	// 取得訂閱對象的 ID
	var edges *cntr.Set[int32]
	returnImmediate := false

	if edges, ok = s.AMgr.Edges[user.Index]; !ok {
		returnImmediate = true
	} else if edges.Length() == 0 {
		returnImmediate = true
	}

	if returnImmediate {
		c.Json(ghttp.StatusOK, ghttp.H{
			"n_post": len(agreement.PostMessages),
			"posts":  agreement.PostMessages,
		})
		s.Http.Send(c)
		return
	}

	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.GetSubscribedPosts
	agreement.Cid = c.GetId()
	var err error
	var sTime time.Time

	if ip.StartTime != "" {
		sTime, err = utils.StringToTime(ip.StartTime)

		if err != nil {
			logger.Warn("Invalied start time fotmat, StartTime: %s", ip.StartTime)
			agreement.StartTime = utils.TimeToTimestamp(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC))
		} else {
			agreement.StartTime = utils.TimeToTimestamp(sTime)
		}
	} else {
		agreement.StartTime = utils.TimeToTimestamp(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC))
	}

	if ip.StopTime != "" {
		sTime, err = utils.StringToTime(ip.StopTime)

		if err != nil {
			logger.Warn("Invalied start time fotmat, StopTime: %s", ip.StopTime)
			agreement.StopTime = utils.TimeToTimestamp(time.Now().UTC())
		} else {
			agreement.StopTime = utils.TimeToTimestamp(sTime)
		}
	} else {
		agreement.StopTime = utils.TimeToTimestamp(time.Now().UTC())
	}

	for edge := range edges.Elements {
		agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
			Index: edge,
		})
	}

	bs, err := agreement.Marshal()
	if err != nil {
		msg := "Failed to marshal agreement"
		logger.Error("%s, err: %+v", msg, err)
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 3,
			"msg": msg,
		})
		s.Http.Send(c)
		return
	}
	td := base.NewTransData()
	td.AddByteArray(bs)
	data := td.FormData()
	err = gos.SendToServer(define.PostMessageServer, &data, int32(len(data)))
	if err != nil {
		msg := "Failed to sned to PostMessage server."
		logger.Error("%s, err: %+v", msg, err)
		c.Json(ghttp.StatusInternalServerError, ghttp.H{
			"ret": 4,
			"msg": msg,
		})
		s.Http.Send(c)
	} else {
		logger.Info("Send define.GetSubscribedPosts request: %+v", agreement)
		s.Http.Finish(c)
	}
}
