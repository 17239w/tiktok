package userinfo

import (
	"errors"
	"tiktok/models"
)

var (
	ErrUserNotExist = errors.New("用户不存在或已注销")
)

// FollowList：关注列表
type FollowList struct {
	UserList []*models.UserInfo `json:"user_list"`
}

// QueryFollowListService：查询关注列表的服务
type QueryFollowListService struct {
	userId int64

	userList []*models.UserInfo //model层返回的用户列表

	*FollowList //关注列表
}

// QueryFollowList：查询关注列表
func QueryFollowList(userId int64) (*FollowList, error) {
	return NewQueryFollowListService(userId).Do()
}

// NewQueryFollowListService：创建QueryFollowListService
func NewQueryFollowListService(userId int64) *QueryFollowListService {
	return &QueryFollowListService{userId: userId}
}

// Do：执行查询关注列表的服务
func (service *QueryFollowListService) Do() (*FollowList, error) {
	var err error
	//1.检查参数
	if err = service.checkNum(); err != nil {
		return nil, err
	}
	//2.准备数据
	if err = service.prepareData(); err != nil {
		return nil, err
	}
	//3.打包数据
	if err = service.packData(); err != nil {
		return nil, err
	}
	//返回数据
	return service.FollowList, nil
}

// 1.checkNum：检查参数
func (service *QueryFollowListService) checkNum() error {
	if !models.NewUserInfoDAO().IsUserExistById(service.userId) {
		return ErrUserNotExist
	}
	return nil
}

// 2.prepareData：准备数据
func (service *QueryFollowListService) prepareData() error {
	var userList []*models.UserInfo
	//调用models层，查询用户的关注列表
	err := models.NewUserInfoDAO().QueryFollowListByUserId(service.userId, &userList)
	if err != nil {
		return err
	}
	//遍历关注列表，将isFollow设为true
	for i, _ := range userList {
		userList[i].IsFollow = true
	}
	service.userList = userList
	return nil
}

// 3.packData：打包数据
func (service *QueryFollowListService) packData() error {
	service.FollowList = &FollowList{UserList: service.userList}
	return nil
}
