package models

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitDB()
	code := m.Run()
	os.Exit(code)
}

// 查找用户关注列表
func TestUserInfoDAO_QueryUserFollowsByUserId(t *testing.T) {
	var userList []*UserInfo
	err := NewUserInfoDAO().QueryFollowListByUserId(1, &userList)
	if err != nil {
		panic(err)
	}
	for _, user := range userList {
		fmt.Printf("%#v\n", *user)
	}
}

// 查找用户粉丝列表
func TestUserInfoDAO_QueryUserFollowersByUserId(t *testing.T) {
	var userList []*UserInfo
	err := NewUserInfoDAO().QueryFollowerListByUserId(1, &userList)
	if err != nil {
		panic(err)
	}
	for _, user := range userList {
		fmt.Printf("%#v\n", *user)
	}
}
