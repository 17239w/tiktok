package models

import "sync"

// UserLogin：用户登录表，和UserInfo属于一对一关系
type UserLogin struct {
	Id         int64  `gorm:"primary_key"`
	UserInfoId int64  // 外键
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:200;notnull"`
}

// UserLoginDAO：用户登录DAO
type UserLoginDAO struct {
}

var (
	userLoginDao *UserLoginDAO
	// 单例模式
	userLoginOnce sync.Once
)

// NewUserLoginDAO：创建UserLoginDAO
func NewUserLoginDao() *UserLoginDAO {
	userLoginOnce.Do(func() {
		userLoginDao = new(UserLoginDAO)
	})
	return userLoginDao
}
