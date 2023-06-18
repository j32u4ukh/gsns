package account

import (
	"fmt"
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

var as *AccountServer
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
	// ==================================================
	// 與 Dba Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	anser, err := gos.Listen(gosDefine.Tcp0, define.AccountPort)
	fmt.Printf("AccountServer | Listen to port %d\n", define.AccountPort)

	if err != nil {
		return errors.Wrapf(err, "Failed to listen port %d", define.AccountPort)
	}

	as = NewAccountServer()
	as.Tcp = anser.(*ans.Tcp0Anser)
	as.Tcp.SetWorkHandler(as.Handler)

	// ==================================================
	// 與 Dba Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	var address string = "127.0.0.1"
	asker, err := gos.Bind(define.DbaServer, address, define.DbaPort, gosDefine.Tcp0, base.OnEventsFunc{
		gosDefine.OnConnected: func(any) {
			logger.Info("完成與 Dba Server 建立 TCP 連線")

			// 請求取得用戶資料
			td := base.NewTransData()
			td.AddByte(define.NormalCommand)
			td.AddUInt16(define.GetUserData)
			data := td.FormData()

			// 將註冊結果回傳主伺服器
			err := gos.SendToServer(define.DbaServer, &data, td.GetLength())

			if err != nil {
				logger.Error("Failed to send to dba %d: %v\nError: %+v", define.DbaServer, data, err)
				return
			}
		},
	})

	if err != nil {
		return errors.Wrapf(err, "Failed to bind address %s:%d", address, define.DbaPort)
	}

	dbaAsker = asker.(*ask.Tcp0Asker)
	dbaAsker.SetWorkHandler(as.DbaHandler)
	logger.Info("DbaServer Asker 伺服器初始化完成")
	logger.Info("伺服器初始化完成")

	fmt.Printf("AccountServer | 伺服器初始化完成\n")

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
		as.Run()

		during = time.Since(start)
		if during < frameTime {
			time.Sleep(frameTime - during)
		}
	}
}
