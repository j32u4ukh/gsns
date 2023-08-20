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

func (s *DbaServer) initTables() error {
	tables := map[string]int{
		"Account":     TidAccount,
		"PostMessage": TidPostMessage,
		"Edge":        TidEdge,
	}
	var tableName string
	var tid int
	var err error
	var result *database.SqlResult
	for tableName, tid = range tables {
		err = s.initTable(tableName, tid)
		if err != nil {
			return errors.Wrapf(err, "初始化表格 %s 時發生錯誤", tableName)
		}
		result, err = s.tables[tid].Creater().Exec()
		if err != nil {
			return errors.Wrapf(err, "建立表格 %s 時發生錯誤", tableName)
		}
		serverLogger.Debug("初始化表格 %s\n%+v\n", tableName, result)
	}
	return nil
}

func (s *DbaServer) initTable(tableName string, tid int) error {
	var tableParams *stmt.TableParam
	var columnParams []*stmt.ColumnParam
	var err error
	tableParams, columnParams, err = plugin.GetProtoParams(fmt.Sprintf("../pb/%s.proto", tableName), dialect.MARIA)
	if err != nil {
		return errors.Wrapf(err, "讀取 %s proto 檔時發生錯誤", tableName)
	}
	s.tables[tid] = gosql.NewTable(tableName, tableParams, columnParams, stmt.ENGINE, stmt.COLLATE, dialect.MARIA)
	s.tables[tid].Init(&gosql.TableConfig{
		Db:               s.db,
		DbName:           s.DbName,
		UseAntiInjection: true,
		InsertFunc:       plugin.InsertProto,
		QueryFunc:        plugin.QueryProto,
		UpdateAnyFunc:    plugin.UpdateProto,
		PtrToDbFunc:      plugin.ProtoToDb,
	})
	return nil
}
