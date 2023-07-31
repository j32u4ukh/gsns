package agrt

import (
	"internal/pbgo"
	"sync"

	"github.com/j32u4ukh/gos/base"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var once sync.Once
var pool sync.Pool

func init() {
	once.Do(func() {
		pool = sync.Pool{
			New: func() any {
				return newAgreement()
			},
		}
	})
}

func GetAgreement() *Agreement {
	return pool.Get().(*Agreement)
}

func PutAgreement(a *Agreement) {
	a.Release()
	pool.Put(a)
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
	agreement.Cid = -1
	agreement.Accounts = agreement.Accounts[:0]
	agreement.Users = agreement.Users[:0]
	agreement.PostMessages = agreement.PostMessages[:0]
	agreement.Edges = agreement.Edges[:0]
}
