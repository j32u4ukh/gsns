package account

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"internal/pbgo"
	"strconv"

	"github.com/pkg/errors"
)

func (m *AccountMgr) AddUser(user *pbgo.User) error {
	if !m.users.ContainKey1(user.Index) {
		err := m.users.Add(user.Index, user.Token, user)
		if err != nil {
			return errors.Wrapf(err, "Failed to add user(%+v)", user)
		}
	}
	return nil
}

func (m *AccountMgr) GetUserByToken(token string) (*pbgo.User, bool) {
	return m.users.GetByKey2(token)
}

// 取得不重複 token
func (m *AccountMgr) getToken() string {
	var token string
	var value uint64
	var err error
	err = binary.Read(rand.Reader, binary.BigEndian, &value)
	if err != nil {
		m.serverLogger.Error("value: %d, err: %+v", value, err)
		return ""
	}
	// 確保 token 唯一
	token = strconv.FormatUint(value, 10)
	for m.users.ContainKey2(token) {
		err = binary.Read(rand.Reader, binary.BigEndian, &value)
		if err != nil {
			m.serverLogger.Error("value: %d, err: %+v", value, err)
			return ""
		}
	}
	m.serverLogger.Info("token: %s", token)
	return token
}

func sha256Func(input string) string {
	// 創建 SHA-256 散列對象
	hash := sha256.New()

	// 將輸入數據寫入散列對象
	hash.Write([]byte(input))

	// 計算散列值
	hashValue := hash.Sum(nil)

	// 將散列值轉換為十六進制字符串
	hashHex := hex.EncodeToString(hashValue)

	return hashHex
}
