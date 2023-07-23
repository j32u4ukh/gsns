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
var logger *glog.Logger

func Init() error {
	logger = glog.GetLogger(0)
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
	logger.Info("Listen to port %d", define.PostMessagePort)

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
			logger.Info("完成與 Dba Server 建立 TCP 連線")
		},
	}, &introduction, &heartbeat)

	if err != nil {
		return errors.Wrapf(err, "Failed to bind address %s:%d", address, define.DbaPort)
	}

	dbaAsker = asker.(*ask.Tcp0Asker)
	dbaAsker.SetWorkHandler(server.DbaHandler)
	logger.Info("DbaServer Asker 伺服器初始化完成")
	logger.Info("伺服器初始化完成")

	// =============================================
	// 開始所有已註冊的監聽
	// =============================================
	gos.StartListen()
	logger.Info("開始所有已註冊的監聽")

	// =============================================
	// 開始所有已註冊的連線
	// =============================================
	err = gos.StartConnect()

	if err != nil {
		return errors.Wrap(err, "與 DbaServer 連線時發生錯誤")
	}

	logger.Info("成功與 DbaServer 連線")
	return nil
}

func Run() {
	var start time.Time
	var during, frameTime time.Duration = 0, 20 * time.Millisecond

	for {
		start = time.Now()

		gos.RunAns()
		gos.RunAsk()
		server.Run()

		during = time.Since(start)
		if during < frameTime {
			time.Sleep(frameTime - during)
		}
	}
}
