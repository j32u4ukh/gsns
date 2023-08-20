package pm

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

var server *PostMessageServer
var dbaAsker *ask.Tcp0Asker
var serverLogger *glog.Logger
var clientLogger *glog.Logger

func Init() error {
	serverLogger = glog.GetLogger(1)
	clientLogger = glog.GetLogger(2)
	err := initGos()
	if err != nil {
		return errors.Wrap(err, "Failed to init gos.")
	}
	return nil
}

func initGos() error {
	td := base.NewTransData()
	agreement := agrt.GetAgreement()
	defer agrt.PutAgreement(agreement)
	agreement.Cmd = define.SystemCommand
	agreement.Service = define.Heartbeat
	bs, _ := agreement.Marshal()
	td.AddByteArray(bs)
	heartbeat := td.FormData()
	agreement.Release()
	td.Clear()

	agreement.Cmd = define.SystemCommand
	agreement.Service = define.Introduction
	agreement.Cipher = "GSNS"
	agreement.Identity = define.PostMessageServer
	bs, _ = agreement.Marshal()
	td.AddByteArray(bs)
	introduction := td.FormData()
	td.Clear()

	// ==================================================
	// 與 Dba Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	anser, err := gos.Listen(gosDefine.Tcp0, define.PostMessagePort)
	serverLogger.Info("Listen to port %d", define.PostMessagePort)

	if err != nil {
		return errors.Wrapf(err, "Failed to listen port %d", define.PostMessagePort)
	}

	server = NewPostMessageServer()
	server.Tcp = anser.(*ans.Tcp0Anser)
	server.Tcp.SetWorkHandler(server.Handler)

	// ==================================================
	// 與 Dba Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	var address string = "127.0.0.1"
	asker, err := gos.Bind(define.DbaServer, address, define.DbaPort, gosDefine.Tcp0, base.OnEventsFunc{
		gosDefine.OnConnected: func(any) {
			serverLogger.Info("完成與 Dba Server 建立 TCP 連線")
		},
	}, &introduction, &heartbeat)

	if err != nil {
		return errors.Wrapf(err, "Failed to bind address %s:%d", address, define.DbaPort)
	}

	dbaAsker = asker.(*ask.Tcp0Asker)
	dbaAsker.SetWorkHandler(server.DbaHandler)
	serverLogger.Info("DbaServer Asker 伺服器初始化完成")
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

	if err != nil {
		return errors.Wrap(err, "與 DbaServer 連線時發生錯誤")
	}

	serverLogger.Info("成功與 DbaServer 連線")
	return nil
}

func Run() {
	gos.SetFrameTime(20 * time.Millisecond)
	gos.Run(nil)
}
