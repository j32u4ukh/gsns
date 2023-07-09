package gsns

import (
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
var logger *glog.Logger

func Init() error {
	logger = glog.GetLogger(0)
	err := initGos()
	if err != nil {
		return errors.Wrap(err, "初始化網路底層時發生錯誤")
	}
	err = initData()
	if err != nil {
		return errors.Wrap(err, "載入數據時發生錯誤")
	}
	return nil
}

func initGos() error {
	ms = newMainServer()

	// ==================================================
	// Http Server: 接受來自客戶端的請求
	// ==================================================
	var port int32 = 1023
	anser, err := gos.Listen(gosDefine.Http, port)
	logger.Info("Listen to port %d", port)

	if err != nil {
		return errors.Wrapf(err, "Failed to listen port %d", port)
	}

	ms.SetHttpAnswer(anser.(*ans.HttpAnser))
	ms.HttpHandler(ms.HttpAnswer.Router)
	logger.Info("Http Anser 伺服器初始化完成")

	// ==================================================
	// 與 Account Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	var address string = "127.0.0.1"
	port = 1021
	askAccount, err := gos.Bind(define.AccountServer, address, 1021, gosDefine.Tcp0, base.OnEventsFunc{
		gosDefine.OnConnected: func(any) {
			logger.Info("成功與 AccountServer 連線")
		},
	})

	if err != nil {
		return errors.Wrapf(err, "Failed to bind address %s:%d", address, port)
	}

	accountAsker = askAccount.(*ask.Tcp0Asker)
	accountAsker.SetWorkHandler(ms.AMgr.WorkHandler)
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
	logger.Info("開始所有已註冊的連線")

	if err != nil {
		return errors.Wrap(err, "與 AccountServer 連線時發生錯誤")
	}
	return nil
}

func initData() error {
	// TODO: 生成向 DBA 取得數據的請求
	// TODO: 生成向 Account 取得數據的請求
	return nil
}

func Run() {
	var start time.Time
	var during, frameTime time.Duration = 0, 20 * time.Millisecond

	for {
		start = time.Now()

		gos.RunAns()
		gos.RunAsk()
		ms.Run()

		during = time.Since(start)
		if during < frameTime {
			time.Sleep(frameTime - during)
		}
	}
}
