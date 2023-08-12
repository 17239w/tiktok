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

func TestUserInfoDAO_QueryUserFollowsByUserId(t *testing.T) {
	var userList []*UserInfo
	err := NewUserInfoDAO().QueryUserFollowsByUserId(1, &userList)
	if err != nil {
		panic(err)
	}
	for _, user := range userList {
		fmt.Printf("%#v\n", *user)
	}
}

func TestUserInfoDAO_QueryUserFollowersByUserId(t *testing.T) {
	var userList []*UserInfo
	err := NewUserInfoDAO().QueryUserFollowersByUserId(1, &userList)
	if err != nil {
		panic(err)
	}
	for _, user := range userList {
		fmt.Printf("%#v\n", *user)
	}
}
