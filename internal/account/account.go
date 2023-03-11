package account

import (
	"fmt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/glog"
	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	"github.com/j32u4ukh/gos/ask"
	"github.com/j32u4ukh/gos/base"
	gosDefine "github.com/j32u4ukh/gos/define"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var as *AccountServer
var dbaAsker *ask.Tcp0Asker
var logger *glog.Logger

func Init(lg *glog.Logger) error {
	logger = lg
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

	as = &AccountServer{}
	as.Tcp = anser.(*ans.Tcp0Anser)
	as.Tcp.SetWorkHandler(as.Handler)
	logger.Info("完成與 Dba Server 建立 TCP 連線")

	// ==================================================
	// 與 Dba Server 建立 TCP 連線，將數據依序寫入緩存
	// ==================================================
	var address string = "127.0.0.1"
	asker, err := gos.Bind(define.DbaServer, address, define.DbaPort, gosDefine.Tcp0)

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

type AccountServer struct {
	Tcp          *ans.Tcp0Anser
	MainServerId int32
}

func (s *AccountServer) Handler(work *base.Work) {
	cmd := work.Body.PopByte()
	logger.Info("cmd: %d", cmd)

	switch cmd {
	case define.SystemCommand:
		s.handleSystemCommand(work)
	case define.CommissionCommand:
		s.handleCommission(work)
	default:
		fmt.Printf("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (rrs *AccountServer) Run() {

}

func (s *AccountServer) handleSystemCommand(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {
	// 回應心跳包
	case define.Heartbeat:
		fmt.Printf("Heart beat! Now: %+v\n", time.Now())
		work.Body.Clear()
		work.Body.AddByte(0)
		work.Body.AddUInt16(0)
		work.Body.AddString("OK")
		work.SendTransData()
	default:
		fmt.Printf("Unsupport service: %d\n", service)
		work.Finish()
	}
}

func (s *AccountServer) handleCommission(work *base.Work) {
	commission := work.Body.PopUInt16()
	logger.Info("commission: %d", commission)

	switch commission {
	case 1023:
		cid := work.Body.PopInt32()
		work.Body.Clear()

		work.Body.AddByte(1)
		work.Body.AddUInt16(1023)
		work.Body.AddInt32(cid)
		work.Body.AddString("Commission completed.")
		work.SendTransData()

	case define.Register:
		s.MainServerId = work.Index
		cid := work.Body.PopInt32()
		bs := work.Body.PopByteArray()
		logger.Info("MainServerId: %d, cid: %d, bs: %+v", s.MainServerId, cid, bs)
		work.Finish()

		td := base.NewTransData()
		td.AddByte(define.CommissionCommand)
		td.AddUInt16(define.Register)
		td.AddInt32(cid)

		// Account data for register
		td.AddByteArray(bs)

		data := td.FormData()
		logger.Info("data: %+v", data)

		// 將註冊數據傳到 Dba 伺服器
		err := gos.SendToServer(define.DbaServer, &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to server %d: %v\nError: %+v", define.DbaServer, data, err)
			return
		}

	default:
		fmt.Printf("Unsupport commission: %d\n", commission)
		work.Finish()
	}
}

func (s *AccountServer) DbaHandler(work *base.Work) {
	cmd := work.Body.PopByte()
	logger.Info("cmd: %d", cmd)

	switch cmd {
	case define.SystemCommand:
		s.handleDbaSystemCommand(work)
	case define.CommissionCommand:
		s.handleDbaCommission(work)
	default:
		fmt.Printf("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (s *AccountServer) handleDbaSystemCommand(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {
	// 回應心跳包
	case define.Heartbeat:
		fmt.Printf("Heart response Now: %+v\n", time.Now())
		work.Finish()
	default:
		fmt.Printf("Unsupport service: %d\n", service)
		work.Finish()
	}
}

func (s *AccountServer) handleDbaCommission(work *base.Work) {
	commission := work.Body.PopUInt16()

	switch commission {
	case define.Register:
		cid := work.Body.PopInt32()
		returnCode := work.Body.PopUInt16()
		bs := work.Body.PopByteArray()
		work.Finish()

		account := &pbgo.Account{}
		err := proto.Unmarshal(bs, account)

		if err != nil {
			return
		}

		logger.Info("New account created : %+v", account)

		td := base.NewTransData()
		td.AddByte(define.CommissionCommand)
		td.AddUInt16(define.Register)
		td.AddInt32(cid)
		td.AddUInt16(returnCode)

		// Account data for register
		td.AddByteArray(bs)

		data := td.FormData()

		// 將註冊數據傳到 Dba 伺服器
		err = gos.SendToClient(define.AccountPort, s.MainServerId, &data, td.GetLength())

		if err != nil {
			logger.Error("Failed to send to client %d: %v\nError: %+v", s.MainServerId, data, err)
			return
		}

	default:
		fmt.Printf("Unsupport commission: %d\n", commission)
		work.Finish()
	}
}
