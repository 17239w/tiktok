package video

import (
	"errors"
	"tiktok/cache"
	"tiktok/models"
)

// List：视频列表
type List struct {
	Videos []*models.Video `json:"video_list,omitempty"`
}

// QueryVideoListByUserIdService：查询用户视频列表的服务
type QueryVideoListByUserIdService struct {
	userId int64
	videos []*models.Video

	videoList *List
}

// QueryVideoListService：查询视频列表服务
func QueryVideoListByUserId(userId int64) (*List, error) {
	return NewQueryVideoListByUserIdService(userId).Do()
}

// NewQueryVideoListByUserIdService：创建查询用户视频列表的服务
func NewQueryVideoListByUserIdService(userId int64) *QueryVideoListByUserIdService {
	return &QueryVideoListByUserIdService{userId: userId}
}

// Do：执行查询用户视频列表的流程
func (service *QueryVideoListByUserIdService) Do() (*List, error) {
	//1.检查参数
	if err := service.checkNum(); err != nil {
		return nil, err
	}
	//2.打包数据
	if err := service.packData(); err != nil {
		return nil, err
	}
	return service.videoList, nil
}

// 1.checkNum：检查参数
func (service *QueryVideoListByUserIdService) checkNum() error {
	//调用models层，检查用户是否存在
	if !models.NewUserInfoDAO().IsUserExistById(service.userId) {
		return errors.New("用户不存在")
	}
	return nil
}

// 2.packData：打包数据
// 注意：Video由于在数据库中没有存储作者信息，所以需要手动填充
func (service *QueryVideoListByUserIdService) packData() error {
	//调用models层，查询视频列表
	err := models.NewVideoDAO().QueryVideoListByUserId(service.userId, &service.videos)
	if err != nil {
		return err
	}
	//调用models层，查询用户信息
	var userInfo models.UserInfo
	err = models.NewUserInfoDAO().QueryUserInfoById(service.userId, &userInfo)
	proxy := cache.NewProxyIndexMap()
	if err != nil {
		return err
	}
	//填充信息(Author和IsFavorite字段)
	for i := range service.videos {
		service.videos[i].Author = userInfo
		service.videos[i].IsFavorite = proxy.GetVideoFavorState(service.userId, service.videos[i].Id)
	}
	service.videoList = &List{Videos: service.videos}
	return nil
}
