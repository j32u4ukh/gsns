package agrt

import (
	"internal/pbgo"
	"sync"

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
			Accounts: []*pbgo.Account{},
			Users:    []*pbgo.User{},
		},
	}
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
		return nil, errors.Wrap(err, "Failed to maeshal Agreement.")
	}
	return bs, nil
}

func (a *Agreement) Release() {}
