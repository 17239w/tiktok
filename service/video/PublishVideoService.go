package video

import (
	"tiktok/models"
	"tiktok/utils"
)

// PublishVideoService：投稿视频服务
type PublishVideoService struct {
	videoName string
	coverName string
	title     string
	userId    int64

	video *models.Video
}

// PostVideo：投稿视频
func PostVideo(userId int64, videoName, coverName, title string) error {
	return NewPublishVideoService(userId, videoName, coverName, title).Do()
}

// NewPublishVideoService：创建PublishVideoService
func NewPublishVideoService(userId int64, videoName, coverName, title string) *PublishVideoService {
	return &PublishVideoService{
		videoName: videoName,
		coverName: coverName,
		userId:    userId,
		title:     title,
	}
}

// Do：执行投稿视频服务
func (service *PublishVideoService) Do() error {
	//1.准备参数
	service.prepareParam()
	//2.组合并添加到数据库
	if err := service.publish(); err != nil {
		return err
	}
	return nil
}

// 1.prepareParam：准备参数
func (service *PublishVideoService) prepareParam() {
	// 获取视频url
	service.videoName = utils.GetFileUrl(service.videoName)
	// 获取封面url
	service.coverName = utils.GetFileUrl(service.coverName)
}

// 2.publish：组合并添加到数据库
func (service *PublishVideoService) publish() error {
	// 组合
	video := &models.Video{
		UserInfoId: service.userId,
		PlayUrl:    service.videoName,
		CoverUrl:   service.coverName,
		Title:      service.title,
	}
	// 调用models层，添加视频到数据库
	return models.NewVideoDAO().AddVideo(video)
}
