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

// [endpoint]/post
// TODO: HTTP 請求處理過程中若失敗，要返回錯誤訊息給客戶端，而非印出日誌或直接返回
func (m *PostMessageMgr) HttpHandler(router *ans.Router) {
	router.POST("/", m.addNewPost)
	router.PATCH("/", m.modifyPost)
	router.GET("/<post_id int>", m.getPost)
	router.GET("/mine", m.getMyPosts)
}

// 用於新增貼文
// [endpoint]/post
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

// 用於讀取特定貼文
// [endpoint]/post/<post_id int>
func (m *PostMessageMgr) getPost(c *ghttp.Context) {
	value := c.GetValue("post_id")
	if value == nil {
		msg := "Failed to get post id."
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		m.httpAnswer.Send(c)
		return
	}

	post_id := value.(int64)
	m.logger.Info("post_id: %d", post_id)

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.GetPost
	agreement.Cid = c.GetId()
	agreement.PostMessages = append(agreement.PostMessages, &pbgo.PostMessage{
		Id: uint64(post_id),
	})
	bs, err := agreement.Marshal()
	if err != nil {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": "Failed to marshal agreement.",
		})
		m.httpAnswer.Send(c)
		return
	}
	td := base.NewTransData()
	td.AddByteArray(bs)
	data := td.FormData()
	err = gos.SendToServer(define.PostMessageServer, &data, int32(len(data)))

	if err != nil {
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 3,
			"msg": "Failed to send data to PostMessage server.",
		})
		m.httpAnswer.Send(c)
	} else {
		m.logger.Info("Send define.GetPost request: %+v", agreement)
		// 將當前 Http 的工作結束
		m.httpAnswer.Finish(c)
	}
}

// 用於讀取貼文
// [endpoint]/post/mine
func (m *PostMessageMgr) getMyPosts(c *ghttp.Context) {
	pmp := &PostMessageProtocol{}
	c.ReadJson(pmp)
	m.logger.Info("PostMessageProtocol: %+v", pmp)
	if pmp.Token == 0 {
		msg := fmt.Sprintf("缺少參數, PostMessage: %+v", pmp)
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 1,
			"msg": msg,
		})
		m.httpAnswer.Send(c)
		return
	}

	user, ok := m.getUserByTokenFunc(pmp.Token)

	if !ok {
		msg := fmt.Sprintf("Not found token %d", pmp.Token)
		m.logger.Error(msg)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": msg,
		})
		m.httpAnswer.Send(c)
		return
	}

	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.CommissionCommand
	agreement.Service = define.GetMyPosts
	agreement.Cid = c.GetId()
	agreement.Accounts = append(agreement.Accounts, &pbgo.Account{
		Index: user.Index,
	})

	// 寫入 agreement
	td := base.NewTransData()
	bs, err := agreement.Marshal()
	if err != nil {
		msg := "Failed to marshal agreement"
		m.logger.Error(fmt.Sprintf("%s, err: %+v", msg, err))
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"ret": 2,
			"msg": msg,
		})
		m.httpAnswer.Send(c)
	}

	td.AddByteArray(bs)
	data := td.FormData()
	m.logger.Info("data: %+v", data)

	// 將數據傳到 PostMessage 伺服器
	err = gos.SendToServer(define.PostMessageServer, &data, int32(len(data)))

	if err != nil {
		m.logger.Error("Failed to send to PostMessage, err: %+v", err)
		c.Json(ghttp.StatusBadRequest, ghttp.H{
			"err": "Failed to send to server.",
		})
		m.httpAnswer.Send(c)
	} else {
		m.logger.Info("Send define.GetMyPosts request: %+v", agreement)
		// 將當前 Http 的工作結束
		m.httpAnswer.Finish(c)
	}
}

// 用於編輯貼文
// [endpoint]/post
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
