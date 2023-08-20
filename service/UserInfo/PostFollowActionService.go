package userinfo

import (
	"errors"
	"tiktok/cache"
	"tiktok/models"
)

const (
	FOLLOW = 1
	CANCEL = 2
)

var (
	ErrIvdAct    = errors.New("未定义操作")
	ErrIvdFolUsr = errors.New("关注用户不存在")
)

// PostFollowActionService：关注/取消关注的流程对象
type PostFollowActionService struct {
	userId     int64
	userToId   int64 //被关注的用户id
	actionType int   //操作类型
}

// PostFollowAction：关注/取消关注操作
func PostFollowAction(userId, userToId int64, actionType int) error {
	return NewPostFollowActionService(userId, userToId, actionType).Do()
}

// NewPostFollowActionService：创建PostFollowActionService
func NewPostFollowActionService(userId int64, userToId int64, actionType int) *PostFollowActionService {
	return &PostFollowActionService{userId: userId, userToId: userToId, actionType: actionType}
}

// Do：执行 关注/取消关注 的操作
func (p *PostFollowActionService) Do() error {
	var err error
	//1.检查参数
	if err = p.checkNum(); err != nil {
		return err
	}
	//2.执行关注/取消关注的操作
	if err = p.publish(); err != nil {
		return err
	}
	return nil
}

// checkNum：检查参数 userId、followId、actionType
func (service *PostFollowActionService) checkNum() error {
	//由于userId经过了token鉴权，只需要检查userToId是否存在
	//调用models层，检查userToId是否存在
	if !models.NewUserInfoDAO().IsUserExistById(service.userToId) {
		return ErrIvdFolUsr
	}
	//检查actionType
	if service.actionType != FOLLOW && service.actionType != CANCEL {
		return ErrIvdAct
	}
	//自己不能关注自己
	if service.userId == service.userToId {
		return ErrIvdAct
	}
	return nil
}

// publish：执行关注/取消关注的操作
func (service *PostFollowActionService) publish() error {
	userDAO := models.NewUserInfoDAO()
	var err error
	switch service.actionType {
	case FOLLOW:
		err = userDAO.AddUserFollow(service.userId, service.userToId)
		//更新redis的关注信息：通过代理对象proxyIndexOperation更新用户关注的状态
		cache.NewProxyIndexMap().UpdateUserRelation(service.userId, service.userToId, true)
	case CANCEL:
		err = userDAO.CancelUserFollow(service.userId, service.userToId)
		//更新redis的关注信息：通过代理对象proxyIndexOperation更新用户关注的状态
		cache.NewProxyIndexMap().UpdateUserRelation(service.userId, service.userToId, false)
	default:
		return ErrIvdAct
	}
	return err
}
