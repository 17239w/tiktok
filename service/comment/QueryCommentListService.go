package comment

import (
	"errors"
	"fmt"
	"tiktok/models"
	"tiktok/utils"
)

// List：评论列表
type List struct {
	Comments []*models.Comment `json:"comment_list"` //models层的Comment
}

// QueryCommentList：查询评论列表
func QueryCommentList(userId, videoId int64) (*List, error) {
	return NewQueryCommentListService(userId, videoId).Do()
}

// QueryCommentListService：查询评论列表的服务
type QueryCommentListService struct {
	userId  int64
	videoId int64

	comments []*models.Comment //models层的Comment

	commentList *List //评论列表
}

// NewQueryCommentListService：创建QueryCommentListService
func NewQueryCommentListService(userId, videoId int64) *QueryCommentListService {
	return &QueryCommentListService{userId: userId, videoId: videoId}
}

// Do：执行查询评论列表的流程
func (service *QueryCommentListService) Do() (*List, error) {
	//1.检查参数
	if err := service.checkNum(); err != nil {
		return nil, err
	}
	//2.准备数据
	if err := service.prepareData(); err != nil {
		return nil, err
	}
	//3.打包数据
	if err := service.packData(); err != nil {
		return nil, err
	}
	return service.commentList, nil
}

// 1.checkNum：检查参数
func (service *QueryCommentListService) checkNum() error {
	//调用models层，通过userId检查用户是否存在
	if !models.NewUserInfoDAO().IsUserExistById(service.userId) {
		return fmt.Errorf("用户%d处于登出状态", service.userId)
	}
	//调用models层，通过videoId检查视频是否存在
	if !models.NewVideoDAO().IsVideoExistById(service.videoId) {
		return fmt.Errorf("视频%d不存在或已经被删除", service.videoId)
	}
	return nil
}

// 2.prepareData：准备数据
func (service *QueryCommentListService) prepareData() error {
	//调用models层，通过VideoId查询评论列表
	err := models.NewCommentDAO().QueryCommentListByVideoId(service.videoId, &service.comments)
	if err != nil {
		return err
	}
	//根据前端的要求填充正确的时间格式
	err = utils.FillCommentListFields(&service.comments)
	if err != nil {
		return errors.New("暂时还没有人评论")
	}
	return nil
}

// 3.打包数据
func (service *QueryCommentListService) packData() error {
	service.commentList = &List{Comments: service.comments}
	return nil
}
