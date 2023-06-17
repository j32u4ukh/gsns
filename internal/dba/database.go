package dba

import (
	"fmt"

	"github.com/j32u4ukh/gosql"
	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/plugin"
	"github.com/j32u4ukh/gosql/stmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/pkg/errors"
)

func (s *DbaServer) initDatabase(db *database.Database, dbName string) {
	s.db = db
	s.DbName = dbName
}

func (s *DbaServer) initTable(tid int) {
	s.tables[tid].Init(&gosql.TableConfig{
		Db:               s.db,
		DbName:           s.DbName,
		UseAntiInjection: false,
		PtrToDbFunc:      plugin.ProtoToDb,
		InsertFunc:       plugin.InsertProto,
		QueryFunc:        plugin.QueryProto,
		UpdateAnyFunc:    plugin.UpdateProto,
	})
}

func (s *DbaServer) initAccountData() error {
	tableName := "Account"
	tableParams, columnParams, err := plugin.GetProtoParams(fmt.Sprintf("../pb/%s.proto", tableName), dialect.MARIA)
	if err != nil {
		return errors.Wrapf(err, "讀取 %s proto 檔時發生錯誤", tableName)
	}
	s.tables[TidAccount] = gosql.NewTable(tableName, tableParams, columnParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	s.initTable(TidAccount)
	result, err := s.tables[TidAccount].Creater().Exec()
	if err != nil {
		return errors.Wrapf(err, "初始化表格 %s 時發生錯誤", tableName)
	}
	logger.Debug("初始化表格 result\n%+v\n", result)
	return nil
}
