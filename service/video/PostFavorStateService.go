package video

import (
	"errors"
	"tiktok/cache"
	"tiktok/models"
)

const (
	PLUS  = 1 //点赞
	MINUS = 2 //取消点赞
)

// PostFavorStateService：点赞或者取消点赞的服务
type PostFavorStateService struct {
	userId     int64
	videoId    int64
	actionType int64 //PLUS or MINUS
}

// PostFavorState：点赞或者取消点赞的操作
func PostFavorState(userId, videoId, actionType int64) error {
	return NewPostFavorStateService(userId, videoId, actionType).Do()
}

// NewPostFavorStateService：创建点赞或者取消点赞的服务
func NewPostFavorStateService(userId, videoId, action int64) *PostFavorStateService {
	return &PostFavorStateService{
		userId:     userId,
		videoId:    videoId,
		actionType: action,
	}
}

// Do：执行点赞或者取消点赞的流程
func (service *PostFavorStateService) Do() error {
	var err error
	//1.检查参数
	if err = service.checkNum(); err != nil {
		return err
	}
	//2.根据actionType执行不同的操作
	switch service.actionType {
	case PLUS:
		err = service.PlusOperation()
	case MINUS:
		err = service.MinusOperation()
	default:
		return errors.New("未定义的操作")
	}
	return err
}

// checkNum：检查参数
func (service *PostFavorStateService) checkNum() error {
	//调用models层，检查用户是否存在
	if !models.NewUserInfoDAO().IsUserExistById(service.userId) {
		return errors.New("用户不存在")
	}
	//检查actionType是否合法
	if service.actionType != PLUS && service.actionType != MINUS {
		return errors.New("未定义的行为")
	}
	return nil
}

// PlusOperation：点赞操作
func (service *PostFavorStateService) PlusOperation() error {
	//调用models层，视频点赞数目+1
	err := models.NewVideoDAO().PlusOneFavorByUserIdAndVideoId(service.userId, service.videoId)
	if err != nil {
		return errors.New("不要重复点赞")
	}
	//调用cache层，更新用户是否点赞的状态
	cache.NewProxyIndexMap().UpdateVideoFavorState(service.userId, service.videoId, true)
	return nil
}

// MinusOperation：取消点赞
func (service *PostFavorStateService) MinusOperation() error {
	//调用models层，视频点赞数目-1
	err := models.NewVideoDAO().MinusOneFavorByUserIdAndVideoId(service.userId, service.videoId)
	if err != nil {
		return errors.New("点赞数目已经为0")
	}
	//调用cache层，更新用户是否点赞的状态
	cache.NewProxyIndexMap().UpdateVideoFavorState(service.userId, service.videoId, false)
	return nil
}
