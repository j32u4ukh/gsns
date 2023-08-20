package main

import (
	"fmt"
	"internal/pbgo"

	"github.com/j32u4ukh/gosql/database"
	"github.com/j32u4ukh/gosql/proto/gstmt"
	"github.com/j32u4ukh/gosql/stmt/dialect"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var gs *gstmt.Gstmt

func main() {
	conf, err := database.NewConfig("../config.yaml")
	if err != nil {
		err = errors.Wrapf(err, "讀取 Config 檔時發生錯誤, err: %+v", err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	dc := conf.GetDatabase()
	gs, err = gstmt.SetGstmt(0, dc.DbName, dialect.MARIA)
	gs.UseAntiInjection(true)
	var sql string
	sql, err = gs.CreateTable(0, "../../cmd/pb", "Account")
	fmt.Printf("CreateTable sql: %s\n", sql)

	account := &pbgo.Account{
		Account:  "acc",
		Password: "password",
	}
	sql, err = gs.Insert(0, []protoreflect.ProtoMessage{account})
	if err != nil {
		err = errors.Wrapf(err, "Insert Account 時發生錯誤, err: %+v", err)
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Insert Account: %s\n", sql)
}
