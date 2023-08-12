package account

import (
	"crypto/rand"
	"encoding/binary"
	"internal/pbgo"

	"github.com/pkg/errors"
)

func (m *AccountMgr) AddUser(user *pbgo.SnsUser) error {
	if !m.users.ContainKey1(user.Index) {
		err := m.users.Add(user.Index, user.Token, user)
		if err != nil {
			return errors.Wrapf(err, "Failed to add user(%+v)", user)
		}
	}
	return nil
}

func (m *AccountMgr) GetUserByToken(token uint64) (*pbgo.SnsUser, bool) {
	return m.users.GetByKey2(token)
}

// 取得不重複 token
func (m *AccountMgr) getToken() uint64 {
	var token uint64
	var err error
	err = binary.Read(rand.Reader, binary.BigEndian, &token)
	if err != nil {
		m.logger.Error("token: %d, err: %+v", token, err)
		return 0
	}
	// 確保 token 唯一
	for m.users.ContainKey2(token) {
		err = binary.Read(rand.Reader, binary.BigEndian, &token)
		if err != nil {
			m.logger.Error("token: %d, err: %+v", token, err)
			return 0
		}
	}
	m.logger.Info("token: %d", token)
	return token
}
