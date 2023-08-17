package video

import (
	"errors"
	"tiktok/models"
)

// FavorList：点赞视频
type FavorList struct {
	Videos []*models.Video `json:"video_list"` //models层的Video对象
}

// QueryFavorVideoListService：查询点赞列表流程
type QueryFavorVideoListService struct {
	userId int64

	videos []*models.Video //models层的Video对象

	videoList *FavorList //点赞列表
}

// QueryFavorVideoList：查询点赞列表
func QueryFavorVideoList(userId int64) (*FavorList, error) {
	return NewQueryFavorVideoListService(userId).Do()
}

// NewQueryFavorVideoListService：初始化查询点赞列表流程
func NewQueryFavorVideoListService(userId int64) *QueryFavorVideoListService {
	return &QueryFavorVideoListService{userId: userId}
}

// Do：执行查询点赞列表流程
func (q *QueryFavorVideoListService) Do() (*FavorList, error) {
	//1.检查参数
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	//2.准备数据
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	//3.打包数据
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.videoList, nil
}

// 1.checkNum：检查参数
func (service *QueryFavorVideoListService) checkNum() error {
	//调用models层，检查用户是否存在
	if !models.NewUserInfoDAO().IsUserExistById(service.userId) {
		return errors.New("用户状态异常")
	}
	return nil
}

// 2.prepareData：准备数据
func (service *QueryFavorVideoListService) prepareData() error {
	//调用models层，查询用户的点赞列表
	err := models.NewVideoDAO().QueryFavorVideoListByUserId(service.userId, &service.videos)
	if err != nil {
		return err
	}
	//填充信息(Author和IsFavorite字段，由于是点赞列表，故所有的都是点赞状态)
	for i := range service.videos {
		//新建一个UserInfo对象
		var userInfo models.UserInfo
		//调用models层，查询作者信息
		err = models.NewUserInfoDAO().QueryUserInfoById(service.videos[i].UserInfoId, &userInfo)
		if err == nil {
			//若查询未出错，则更新作者信息；否则不更新作者信息
			service.videos[i].Author = userInfo
		}
		service.videos[i].IsFavorite = true
	}
	return nil
}

// 3.packData：打包数据
func (service *QueryFavorVideoListService) packData() error {
	//将视频列表打包成FavorList
	service.videoList = &FavorList{Videos: service.videos}
	return nil
}
