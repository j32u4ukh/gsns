package account

import "internal/pbgo"

func (m *AccountMgr) AddUser(user *pbgo.SnsUser) {
	m.users[user.Index] = user
}
