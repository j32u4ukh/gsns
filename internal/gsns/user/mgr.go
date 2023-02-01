package user

import (
	"internal/pbgo"

	"github.com/j32u4ukh/glog"
)

type UserMgr struct {
	// TODO: cntr.Array
	users  []*pbgo.SnsUser
	logger *glog.Logger
}

func NewUserMgr(nUser int32, lg *glog.Logger) *UserMgr {
	m := &UserMgr{
		users:  make([]*pbgo.SnsUser, 0, nUser),
		logger: lg,
	}
	return m
}

func (m *UserMgr) AddUser(user *pbgo.SnsUser) error {
	m.users = append(m.users, user)
	return nil
}
