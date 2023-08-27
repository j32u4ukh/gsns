package gsns

import (
	"internal/agrt"
	"internal/define"
	"time"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/ask"
	"github.com/j32u4ukh/gos/base"
	gosDefine "github.com/j32u4ukh/gos/define"
	"github.com/pkg/errors"
)

var ms *MainServer
var accountAsker *ask.Tcp0Asker
var pmAsker *ask.Tcp0Asker
var serverLogger *glog.Logger
var clientLogger *glog.Logger

func Init() error {
	serverLogger = glog.GetLogger(1)
	clientLogger = glog.GetLogger(2)
	err := initGos()
	if err != nil {
		return errors.Wrap(err, "初始化網路底層時發生錯誤")
	}
	return nil
}

func initGos() error {
	ms = newMainServer()
	td := base.NewTransData()
	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.SystemCommand
	agreement.Service = define.Heartbeat
	bs, _ := agreement.Marshal()
	td.AddByteArray(bs)
	data := td.FormData()
	agreement.Release()
	td.Clear()
	heartbeat := make([]byte, len(data))
	copy(heartbeat, data)

	agreement.Cmd = define.SystemCommand
	agreement.Service = define.Introduction
	agreement.Cipher = "GSNS"
	agreement.Identity = define.GsnsServer
	bs, _ = agreement.Marshal()
	td.AddByteArray(bs)
	data = td.FormData()
	td.Clear()
	introduction := make([]byte, len(data))
	copy(introduction, data)

	// ==================================================
	// Http Server: 接受來自客戶端的請求
	// ==================================================
	var port int32 = 1023
	anser, err := gos.Listen(gosDefine.Http, port)
	serverLogger.Info("Listen to port %d", port)

	if err != nil {
		return errors.Wrapf(err, "Failed to listen port %d", port)
	}

	httpServer := anser.(*ans.HttpAnser)
	httpServer.Cors("http://localhost:8080")
	ms.SetHttpAnswer(httpServer)
	serverLogger.Info("Http Anser 伺服器初始化完成")

	// ==================================================
	// 與 Account Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	var address string = "127.0.0.1"
	port = 1021
	askAccount, err := gos.Bind(define.AccountServer, address, 1021, gosDefine.Tcp0, base.OnEventsFunc{
		gosDefine.OnConnected: func(any) {
			serverLogger.Info("成功與 AccountServer 連線")
		},
	}, &introduction, &heartbeat)

	if err != nil {
		return errors.Wrapf(err, "Failed to bind address %s:%d", address, port)
	}

	accountAsker = askAccount.(*ask.Tcp0Asker)
	accountAsker.SetWorkHandler(ms.AMgr.WorkHandler)

	// ==================================================
	// 與 PostMessage Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	askPostMessage, err := gos.Bind(define.PostMessageServer, address, define.PostMessagePort, gosDefine.Tcp0, base.OnEventsFunc{
		gosDefine.OnConnected: func(any) {
			serverLogger.Info("成功與 PostMessage Server 連線")
		},
	}, &introduction, &heartbeat)

	if err != nil {
		return errors.Wrapf(err, "Failed to bind address %s:%d", address, define.PostMessagePort)
	}

	pmAsker = askPostMessage.(*ask.Tcp0Asker)
	pmAsker.SetWorkHandler(ms.PMgr.WorkHandler)

	// =============================================
	serverLogger.Info("伺服器初始化完成")

	// =============================================
	// 開始所有已註冊的監聽
	// =============================================
	gos.StartListen()
	serverLogger.Info("開始所有已註冊的監聽")

	// =============================================
	// 開始所有已註冊的連線
	// =============================================
	err = gos.StartConnect()
	serverLogger.Info("開始所有已註冊的連線")

	if err != nil {
		return errors.Wrap(err, "與 AccountServer 連線時發生錯誤")
	}
	return nil
}

func Run() {
	gos.SetFrameTime(20 * time.Millisecond)
	gos.Run(nil)
}
