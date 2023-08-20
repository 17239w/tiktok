package userinfo

import (
	"tiktok/cache"
	"tiktok/models"
)

// FollowerList：粉丝列表
type FollowerList struct {
	UserList []*models.UserInfo `json:"user_list"`
}

// QueryFollowerListService：查询粉丝列表的服务
type QueryFollowerListService struct {
	userId int64

	userList []*models.UserInfo

	*FollowerList
}

// QueryFollowerList：查询粉丝列表
func QueryFollowerList(userId int64) (*FollowerList, error) {
	return NewQueryFollowerListService(userId).Do()
}

// NewQueryFollowerListService：创建QueryFollowerListService
func NewQueryFollowerListService(userId int64) *QueryFollowerListService {
	return &QueryFollowerListService{userId: userId}
}

// Do：执行查询粉丝列表的流程
func (service *QueryFollowerListService) Do() (*FollowerList, error) {
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
	return service.FollowerList, nil
}

// 1.checkNum：检查参数
func (service *QueryFollowerListService) checkNum() error {
	//调用models层，判断用户是否存在
	if !models.NewUserInfoDAO().IsUserExistById(service.userId) {
		return ErrUserNotExist
	}
	return nil
}

// 2.prepareData：准备数据
func (service *QueryFollowerListService) prepareData() error {
	//调用models层，查询用户的粉丝列表
	err := models.NewUserInfoDAO().QueryFollowerListByUserId(service.userId, &service.userList)
	if err != nil {
		return err
	}
	//填充is_follow字段
	for _, v := range service.userList {
		//调用cache层，查询是否关注
		v.IsFollow = cache.NewProxyIndexMap().GetUserRelation(service.userId, v.Id)
	}
	return nil
}

// 3.packData：打包数据
func (service *QueryFollowerListService) packData() error {
	service.FollowerList = &FollowerList{UserList: service.userList}

	return nil
}
