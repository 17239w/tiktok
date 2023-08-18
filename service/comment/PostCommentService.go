package comment

import (
	"errors"
	"fmt"
	"tiktok/models"
	"tiktok/utils"
)

const (
	CREATE = 1
	DELETE = 2
)

// Response：评论的响应
type Response struct {
	MyComment *models.Comment `json:"comment"` //评论信息
}

// PostComment：增加评论
func PostComment(userId int64, videoId int64, commentId int64, actionType int64, commentText string) (*Response, error) {
	return NewPostCommentService(userId, videoId, commentId, actionType, commentText).Do()
}

// PostCommentService：增加评论的服务
type PostCommentService struct {
	userId      int64
	videoId     int64
	commentId   int64
	actionType  int64 //CREATE(1) or DELETE(2)
	commentText string

	comment *models.Comment //models层的comment

	*Response //评论的响应
}

// NewPostCommentService：创建PostCommentService
func NewPostCommentService(userId int64, videoId int64, commentId int64, actionType int64, commentText string) *PostCommentService {
	return &PostCommentService{userId: userId, videoId: videoId, commentId: commentId, actionType: actionType, commentText: commentText}
}

// Do：执行增加评论的流程
func (service *PostCommentService) Do() (*Response, error) {
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
	return service.Response, err
}

// CreateComment：增加评论
func (service *PostCommentService) CreateComment() (*models.Comment, error) {
	comment := models.Comment{UserInfoId: service.userId, VideoId: service.videoId, Content: service.commentText}
	//调用models层，增加comment
	err := models.NewCommentDAO().AddCommentAndUpdateCount(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// DeleteComment：删除评论
func (service *PostCommentService) DeleteComment() (*models.Comment, error) {
	//获取comment
	var comment models.Comment
	//调用models层，通过id查询comment
	err := models.NewCommentDAO().QueryCommentById(service.commentId, &comment)
	if err != nil {
		return nil, err
	}
	//调用models层，删除comment
	err = models.NewCommentDAO().DeleteCommentAndUpdateCountById(service.commentId, service.videoId)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// 1.checkNum：检查参数
func (service *PostCommentService) checkNum() error {
	//调用models层，检查用户和视频是否存在、actionType是否合法
	if !models.NewUserInfoDAO().IsUserExistById(service.userId) {
		return fmt.Errorf("用户%d不存在", service.userId)
	}
	if !models.NewVideoDAO().IsVideoExistById(service.videoId) {
		return fmt.Errorf("视频%d不存在", service.videoId)
	}
	if service.actionType != CREATE && service.actionType != DELETE {
		return errors.New("未定义的行为")
	}
	return nil
}

// 2.prepareData：准备数据
func (service *PostCommentService) prepareData() error {
	var err error
	switch service.actionType {
	case CREATE:
		//调用service层，增加评论
		service.comment, err = service.CreateComment()
	case DELETE:
		//调用service层，删除评论
		service.comment, err = service.DeleteComment()
	default:
		return errors.New("未定义的操作")
	}
	return err
}

// 3.packData：打包数据
func (service *PostCommentService) packData() error {
	userInfo := models.UserInfo{}
	//调用models层，通过id查询userInfo
	_ = models.NewUserInfoDAO().QueryUserInfoById(service.comment.UserInfoId, &userInfo)
	service.comment.User = userInfo
	//调用utils层，填充评论
	_ = utils.FillCommentFields(service.comment)
	service.Response = &Response{MyComment: service.comment}

	return nil
}
