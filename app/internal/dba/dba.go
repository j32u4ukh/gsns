package dba

import (
	"internal/define"
	"time"

	"github.com/j32u4ukh/glog/v2"
	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	gosDefine "github.com/j32u4ukh/gos/define"
	"github.com/j32u4ukh/gosql/database"
	"github.com/pkg/errors"
)

var s *DbaServer
var serverLogger *glog.Logger
var clientLogger *glog.Logger

func Init() error {
	serverLogger = glog.GetLogger(1)
	clientLogger = glog.GetLogger(2)
	err := initGos()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize gos.")
	}
	err = initDatabase()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize database.")
	}
	return nil
}

// 初始化伺服器連線與監聽
func initGos() error {
	anser, err := gos.Listen(gosDefine.Tcp0, define.DbaPort)
	serverLogger.Info("Listen to port %d", define.DbaPort)

	if err != nil {
		return errors.Wrapf(err, "Failed to listen to port %d.", define.DbaPort)
	}

	s = NewDbaServer()
	dbaAnser := anser.(*ans.Tcp0Anser)
	dbaAnser.SetWorkHandler(s.Handler)
	serverLogger.Info("伺服器初始化完成")

	// =============================================
	// 開始所有已註冊的監聽
	// =============================================
	gos.StartListen()
	serverLogger.Info("開始所有已註冊的監聽")
	return nil
}

// 初始化資料庫相關
func initDatabase() error {
	conf, err := database.NewConfig("./config.yaml")
	if err != nil {
		return errors.Wrapf(err, "讀取 Config 檔時發生錯誤, err: %+v", err)
	}

	dc := conf.GetDatabase()
	db, err := database.Connect(0, dc.User, dc.Password, dc.Host, dc.Port, dc.DbName)

	if err != nil {
		return errors.Wrapf(err, "與資料庫連線時發生錯誤, err: %+v", err)
	}

	s.initDatabase(db, dc.DbName)

	if err != nil {
		return errors.Wrapf(err, "Failed to init database, err: %+v", err)
	}

	err = s.initTables()

	if err != nil {
		return errors.Wrapf(err, "Failed to init tables, err: %+v", err)
	}
	return nil
}

func Run() {
	defer s.db.Close()
	gos.SetFrameTime(20 * time.Millisecond)
	gos.Run(nil)
}
