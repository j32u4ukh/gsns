package agrt

import (
	"internal/define"
	"internal/pbgo"
	"sync"

	"github.com/j32u4ukh/gos"
	"github.com/j32u4ukh/gos/base"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var once sync.Once
var agrementPool sync.Pool
var transdataPool sync.Pool

func init() {
	once.Do(func() {
		agrementPool = sync.Pool{
			New: func() any {
				return newAgreement()
			},
		}
		transdataPool = sync.Pool{
			New: func() any {
				return base.NewTransData()
			},
		}
	})
}

func SendWork(work *base.Work, agreement *Agreement) ([]byte, error) {
	bs, err := agreement.Marshal()
	if err != nil {
		work.Finish()
		return nil, errors.Wrap(err, "Failed to marshal agreement")
	}
	work.Body.AddByteArray(bs)
	work.SendTransData()
	return bs, nil
}

func SendToServer(serverId int32, agreement *Agreement) ([]byte, error) {
	// 寫入 agreement
	td := transdataPool.Get().(*base.TransData)
	defer func() {
		td.Clear()
		transdataPool.Put(td)
	}()
	bs, err := agreement.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal agreement")
	}
	td.AddByteArray(bs)
	data := td.FormData()

	// 將註冊數據傳到伺服器
	err = gos.SendToServer(serverId, &data, int32(len(data)))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to sned to %s", define.ServerName(serverId))
	}
	return data, nil
}

func SendToClient(serverId int32, cid int32, agreement *Agreement) ([]byte, error) {
	// 寫入 agreement
	td := transdataPool.Get().(*base.TransData)
	defer func() {
		td.Clear()
		transdataPool.Put(td)
	}()
	bs, err := agreement.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal agreement")
	}
	td.AddByteArray(bs)
	data := td.FormData()

	// 將註冊數據傳到伺服器
	err = gos.SendToClient(serverId, cid, &data, int32(len(data)))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to sned to %s", define.ServerName(serverId))
	}
	return data, nil
}

func GetAgreement() *Agreement {
	return agrementPool.Get().(*Agreement)
}

func PutAgreement(a *Agreement) {
	a.Release()
	agrementPool.Put(a)
}

type Agreement struct {
	*pbgo.Agreement
}

func newAgreement() *Agreement {
	return &Agreement{
		Agreement: &pbgo.Agreement{
			Accounts:     []*pbgo.Account{},
			Users:        []*pbgo.User{},
			PostMessages: []*pbgo.PostMessage{},
			Edges:        []*pbgo.Edge{},
		},
	}
}

func (a *Agreement) Init(work *base.Work) error {
	defer work.Body.Clear()
	bs := work.Body.PopByteArray()
	err := a.Unmarshal(bs)
	if err != nil {
		return errors.Wrap(err, "Failed to init Agreement.")
	}
	return nil
}

func (a *Agreement) Unmarshal(bs []byte) error {
	err := proto.Unmarshal(bs, a.Agreement)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal Agreement.")
	}
	return nil
}

func (a *Agreement) Marshal() ([]byte, error) {
	bs, err := proto.Marshal(a.Agreement)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal Agreement.")
	}
	return bs, nil
}

func (a *Agreement) Release() {
	agreement := a.Agreement
	agreement.Cmd = -1
	agreement.Service = -1
	agreement.ReturnCode = -1
	agreement.Msg = ""
	agreement.Cid = 0
	agreement.Accounts = agreement.Accounts[:0]
	agreement.Users = agreement.Users[:0]
	agreement.PostMessages = agreement.PostMessages[:0]
	agreement.Cipher = ""
	agreement.Identity = 0
	agreement.Edges = agreement.Edges[:0]
	agreement.StartUtc = 0
	agreement.StopUtc = 0
}
