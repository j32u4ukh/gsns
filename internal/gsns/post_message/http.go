package pm

import (
	"fmt"
	"internal/agrt"
	"internal/define"
	"internal/pbgo"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gos/base/ghttp"
)

// TODO: HTTP 請求處理過程中若失敗，要返回錯誤訊息給客戶端，而非印出日誌或直接返回
func (m *PostMessageMgr) HttpHandler(router *ans.Router) {
	router.POST("/", m.addNewPost)
	router.PATCH("/", m.modifyPost)
	router.GET("/<post_id int>", m.getPost)
}

// 用於新增貼文
func (m *PostMessageMgr) addNewPost(c *ghttp.Context) {
	pmp := &PostMessageProtocol{}
	c.ReadJson(pmp)
	m.logger.Info("PostMessageProtocol: %+v", pmp)

	if pmp.Token == 0 || pmp.Content == "" {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": fmt.Sprintf("缺少參數, PostMessage: %+v", pmp),
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

	// 將註冊數據傳到 PostMessage 伺服器
	err := gos.SendToServer(define.PostMessageServer, &data, int32(len(data)))

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

// 用於讀取貼文
func (m *PostMessageMgr) getPost(c *ghttp.Context) {
	// TODO: 返回指定的貼文內容
	value := c.GetValue("post_id")
	if value != nil {
		post_id := value.(int64)
		m.logger.Info("post_id: %d", post_id)

		agreement := agrt.GetAgreement()
		defer agrt.PutAgreement(agreement)
		agreement.Cmd = define.CommissionCommand
		agreement.Service = define.GetPost
		agreement.Cid = c.GetId()
		pm := &pbgo.PostMessage{
			Id: uint64(post_id),
		}
		agreement.PostMessages = append(agreement.PostMessages, pm)
		bs, _ := agreement.Marshal()
		td := base.NewTransData()
		td.AddByteArray(bs)
		data := td.FormData()
		err := gos.SendToServer(define.PostMessageServer, &data, int32(len(data)))

		if err != nil {
			c.Json(ghttp.StatusBadRequest, ghttp.H{
				"error": "Failed to send data to PostMessage server.",
			})
			m.httpAnswer.Send(c)
		} else {
			// 將當前 Http 的工作結束
			m.httpAnswer.Finish(c)
		}

	} else {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"error": "Failed to get post id.",
		})
		m.httpAnswer.Send(c)
	}
}

// 用於編輯貼文
func (m *PostMessageMgr) modifyPost(c *ghttp.Context) {
	pmp := &PostMessageProtocol{}
	c.ReadJson(pmp)
	m.logger.Info("PostMessageProtocol: %+v", pmp)
	if pmp.Token == 0 || pmp.PostId == 0 || pmp.Content == "" {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": fmt.Sprintf("缺少參數, PostMessage: %+v", pmp),
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
	agreement.Service = define.ModifyPost
	agreement.Cid = c.GetId()
	pm := &pbgo.PostMessage{
		Id:       pmp.PostId,
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

	// 將註冊數據傳到 PostMessage 伺服器
	err := gos.SendToServer(define.PostMessageServer, &data, int32(len(data)))

	if err != nil {
		m.logger.Error("Failed to send to server %d: %v\nError: %+v", define.PostMessageServer, data, err)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": "Failed to send to server.",
		})
		m.httpAnswer.Send(c)
		return
	}

	// 將當前 Http 的工作結束
	m.httpAnswer.Finish(c)
}
