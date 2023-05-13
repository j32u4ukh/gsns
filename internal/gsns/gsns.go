package gsns

import (
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/glog"
	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/ask"
	"github.com/j32u4ukh/gos/define"
	"github.com/pkg/errors"
)

const EDbaServer int32 = 0
const EAccountServer int32 = 1

var ms *MainServer
var dbaAsker *ask.Tcp0Asker
var accountAsker *ask.Tcp0Asker
var logger *glog.Logger

func Init(lg *glog.Logger) error {
	logger = lg
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
	ms = newMainServer(100)

	// ==================================================
	// Http Server: 接受來自客戶端的請求
	// ==================================================
	var port int32 = 1023
	anser, err := gos.Listen(define.Http, port)
	logger.Info("Listen to port %d", port)

	if err != nil {
		return errors.Wrapf(err, "Failed to listen port %d", port)
	}

	ms.HttpAnswer = anser.(*ans.HttpAnser)
	ms.HttpHandler(ms.HttpAnswer.Router)
	logger.Info("Http Anser 伺服器初始化完成")

	// ==================================================
	// 與 Dba Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	var address string = "127.0.0.1"
	port = 1022
	asker, err := gos.Bind(EDbaServer, address, 1022, define.Tcp0, nil)

	if err != nil {
		return errors.Wrapf(err, "Failed to bind address %s:%d", address, port)
	}

	dbaAsker = asker.(*ask.Tcp0Asker)
	dbaAsker.SetWorkHandler(ms.DbaHandler)
	logger.Info("DbaServer Asker 伺服器初始化完成")
	logger.Info("伺服器初始化完成")

	// ==================================================
	// 與 Account Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	port = 1021
	askAccount, err := gos.Bind(EAccountServer, address, 1021, define.Tcp0, nil)

	if err != nil {
		return errors.Wrapf(err, "Failed to bind address %s:%d", address, port)
	}

	accountAsker = askAccount.(*ask.Tcp0Asker)
	accountAsker.SetWorkHandler(ms.AccountHandler)
	logger.Info("AccountServer Asker 伺服器初始化完成")
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

func initData() error {
	ms.AddUser(&pbgo.SnsUser{
		Index: 0,
		Name:  "Henry",
	})
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
