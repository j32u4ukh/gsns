package pm

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
func (m *PostMessageMgr) HttpHandler(router *ans.Router) {
	router.GET("/", func(c *ghttp.Context) {
		c.Json(ghttp.StatusOK, ghttp.H{
			"msg": "Hello post message",
			"cid": c.GetId(),
		})
		m.httpAnswer.Send(c)
	})
	router.POST("/", m.addNewPost)
	router.GET("/<post_id int>", func(c *ghttp.Context) {
		// TODO: 返回指定的貼文內容
		ok, value := c.GetParam("post_id")
		if ok {
			post_id, _ := strconv.Atoi(value)
			m.logger.Info("post_id: %d", post_id)
			c.Json(ghttp.StatusOK, ghttp.H{
				"ret": 0,
				"msg": fmt.Sprintf("post_id: %d", post_id),
			})
		} else {
			c.Json(ghttp.StatusBadRequest, ghttp.H{
				"ret": 1,
				"msg": "Failed to get post id.",
			})
		}
		m.httpAnswer.Send(c)
	})
}

// 用於新增貼文
func (m *PostMessageMgr) addNewPost(c *ghttp.Context) {
	pmp := &PostMessageProtocol{}
	c.ReadJson(pmp)
	m.logger.Info("PostMessageProtocol: %+v", pmp)

	if pmp.Token == 0 || pmp.Content == "" {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": fmt.Sprintf("缺少參數, account: %+v", pmp),
		})
		m.httpAnswer.Send(c)
		return
	}
	user, ok := m.getUserByTokenFunc(pmp.Token)

	if !ok {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"msg": fmt.Sprintf("Not found token %d", pmp.Token),
		})
		m.httpAnswer.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.AddPost
	agreement.Cid = c.GetId()
	pm := &pbgo.PostMessage{
		ParentId: pmp.ParentId,
		UserId:   user.Index,
		Content:  pmp.Content,
	}
	m.logger.Info("PostMessage: %+v", pm)
	agreement.PostMessages = append(agreement.PostMessages, pm)

	// 寫入 agreement
	td := base.NewTransData()
	bs, _ := agreement.Marshal()
	td.AddByteArray(bs)
	data := td.FormData()
	m.logger.Info("data: %+v", data)

	// 將註冊數據傳到 Account 伺服器
	err := gos.SendToServer(define.PostMessageServer, &data, td.GetLength())

	if err != nil {
		m.logger.Error("Failed to send to server %d: %v\nError: %+v", define.DbaServer, data, err)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": "Failed to send to server.",
		})
		m.httpAnswer.Send(c)
		return
	}

	// 將當前 Http 的工作結束
	m.httpAnswer.Finish(c)
}
