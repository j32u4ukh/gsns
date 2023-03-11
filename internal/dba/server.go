package dba

import (
	"fmt"
	"internal/define"
	"internal/pbgo"
	"time"

	"github.com/j32u4ukh/gos/base"
	"github.com/j32u4ukh/gosql/database"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type DbaServer struct {
}

func (s *DbaServer) Handler(work *base.Work) {
	cmd := work.Body.PopByte()

	switch cmd {
	case define.SystemCommand:
		s.handleSystemCommand(work)
	case 1:
	case define.CommissionCommand:
		s.handleCommission(work)

	default:
		fmt.Printf("Unsupport command: %d\n", cmd)
		work.Finish()
	}
}

func (rrs *DbaServer) Run() {

}

func (s *DbaServer) handleSystemCommand(work *base.Work) {
	service := work.Body.PopUInt16()

	switch service {
	// 回應心跳包
	case define.Heartbeat:
		fmt.Printf("Heart beat! Now: %+v\n", time.Now())
		work.Body.Clear()
		work.Body.AddByte(define.SystemCommand)
		work.Body.AddUInt16(define.Heartbeat)
		work.Body.AddString("OK")
		work.SendTransData()
	default:
		fmt.Printf("Unsupport service: %d\n", service)
		work.Finish()
	}
}

func (s *DbaServer) handleCommission(work *base.Work) {
	commission := work.Body.PopUInt16()
	var cid int32 = work.Body.PopInt32()
	logger.Info("commission: %d, cid: %d", commission, cid)

	switch commission {
	case 1023:
		work.Body.Clear()
		work.Body.AddByte(1)
		work.Body.AddUInt16(1023)
		work.Body.AddInt32(cid)
		work.Body.AddString("Commission completed.")
		work.SendTransData()

	case define.Register:
		// 建立使用者資料
		bs := work.Body.PopByteArray()
		account := &pbgo.Account{}
		err := proto.Unmarshal(bs, account)
		work.Body.Clear()
		work.Body.AddByte(define.CommissionCommand)
		work.Body.AddUInt16(define.Register)
		work.Body.AddInt32(cid)

		if err != nil {
			logger.Error("Unmarshal account err: %+v", err)
			// TODO: send error message back to client
		} else {
			var sql string
			sql, err = gs.Insert(TidAccount, []protoreflect.ProtoMessage{account})

			if err != nil {
				fmt.Printf("Error: %+v", err)
				return
			}

			var result *database.SqlResult
			result, err = db.Exec(sql)

			if err != nil {
				fmt.Printf("Insert Exec err: %+v\n", err)
				return
			}

			logger.Info("result: %s", result)
			work.Body.AddUInt16(0)
			account.Index = int32(result.LastInsertId)
			account.Account = ""
			account.Password = ""
			bs, _ = proto.Marshal(account)
			work.Body.AddByteArray(bs)
			work.SendTransData()
		}
	default:
		fmt.Printf("Unsupport commission: %d\n", commission)
		work.Finish()
	}
}
