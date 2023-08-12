package models

import (
	"errors"
	"sync"
)

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

// QueryUserLogin：根据用户名和密码查询用户登录信息
func (dao *UserLoginDAO) QueryUserLogin(username, password string, login *UserLogin) error {
	if login == nil {
		return errors.New("UserLogin is nil")
	}
	DB.Where("username = ? AND password = ?", username, password).First(login)
	if login.Id == 0 {
		return errors.New("user not exist")
	}
	return nil
}

// IsUserExist：根据姓名判断用户是否存在
func (dao *UserLoginDAO) IsUserExist(username string) bool {
	var userlogin UserLogin
	if error := DB.Where("username = ?", username).First(&userlogin).Error; error == nil {
		return true
	}
	return false
}
