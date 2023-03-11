package dba

import (
	"internal/define"
	"time"

	"github.com/j32u4ukh/glog"
	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/ans"
	gosDefine "github.com/j32u4ukh/gos/define"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/proto/gstmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/pkg/errors"
)

var s *DbaServer
var db *database.Database
var gs *gstmt.Gstmt
var logger *glog.Logger

func Init(lg *glog.Logger) error {
	logger = lg
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
	logger.Info("Listen to port %d", define.DbaPort)

	if err != nil {
		return errors.Wrapf(err, "Failed to listen to port %d.", define.DbaPort)
	}

	s = &DbaServer{}
	dbaAnser := anser.(*ans.Tcp0Anser)
	dbaAnser.SetWorkHandler(s.Handler)
	logger.Info("伺服器初始化完成")

	// =============================================
	// 開始所有已註冊的監聽
	// =============================================
	gos.StartListen()
	logger.Info("開始所有已註冊的監聽")
	return nil
}

// 初始化資料庫相關
func initDatabase() error {
	conf, err := database.NewConfig("./config.yaml")
	if err != nil {
		return errors.Wrapf(err, "讀取 Config 檔時發生錯誤, err: %+v", err)
	}

	dc := conf.GetDatabase()
	db, err = database.Connect(0, dc.UserName, dc.Password, dc.Server, dc.Port, dc.Name)

	if err != nil {
		return errors.Wrapf(err, "與資料庫連線時發生錯誤, err: %+v", err)
	}

	gs, err = gstmt.SetGstmt(0, dc.Name, dialect.MARIA)

	if err != nil {
		return errors.Wrapf(err, "SetGstmt err: %+v", err)
	}

	var sql string
	sql, err = gs.CreateTable(TidAccount, "../pb", "Account")
	if err != nil {
		return errors.Wrapf(err, "Create err: %+v", err)
	}

	var result *database.SqlResult
	result, err = db.Exec(sql)

	if err != nil {
		return errors.Wrapf(err, "Create Exec err: %+v", err)
	}

	logger.Debug("result: %s", result)
	gs.UseAntiInjection(true)
	return nil
}

func Run() {
	var start time.Time
	var during, frameTime time.Duration = 0, 20 * time.Millisecond
	defer db.Close()

	for {
		start = time.Now()

		gos.RunAns()
		s.Run()

		during = time.Since(start)
		if during < frameTime {
			time.Sleep(frameTime - during)
		}
	}
}
