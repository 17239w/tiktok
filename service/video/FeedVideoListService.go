package video

import (
	"tiktok/models"
	"tiktok/utils"
	"time"
)

const (
	MaxVideoNum = 30 //每次最多返回的视频流数量
)

// FeedVideoList：视频流列表
type FeedVideoList struct {
	Videos   []*models.Video `json:"video_list,omitempty"` //model层的video
	NextTime int64           `json:"next_time,omitempty"`  //下次请求的时间戳
}

// QueryFeedVideoListService：查询视频流列表的服务
type QueryFeedVideoListService struct {
	userId     int64
	latestTime time.Time

	videos   []*models.Video
	nextTime int64

	feedVideo *FeedVideoList
}

// QueryFeedVideoList：查询视频流列表
func QueryFeedVideoList(userId int64, latestTime time.Time) (*FeedVideoList, error) {
	return NewQueryFeedVideoListService(userId, latestTime).Do()
}

// NewQueryFeedVideoListService：创建QueryFeedVideoListService
func NewQueryFeedVideoListService(userId int64, latestTime time.Time) *QueryFeedVideoListService {
	return &QueryFeedVideoListService{userId: userId, latestTime: latestTime}
}

// Do：执行查询视频流列表的服务
func (service *QueryFeedVideoListService) Do() (*FeedVideoList, error) {
	//1.检查参数
	service.checkNum()
	//2.准备数据
	if err := service.prepareData(); err != nil {
		return nil, err
	}
	//3.打包数据
	if err := service.packData(); err != nil {
		return nil, err
	}
	return service.feedVideo, nil
}

// 1.checkNum：检查参数
func (service *QueryFeedVideoListService) checkNum() {
	//上层通过把userId置零，表示userId不存在或不需要
	//这里说明userId是有效的，可以定制性的做一些登录用户的专属视频推荐
	// if service.userId > 0 {
	// }
	//如果latestTime为零值，说明latestTime是无效的，需要重新赋值
	if service.latestTime.IsZero() {
		service.latestTime = time.Now()
	}
}

// 2.prepareData：准备数据
func (service *QueryFeedVideoListService) prepareData() error {
	//调用model层，查询视频列表
	err := models.NewVideoDAO().GetVideoListByRecentUpload(MaxVideoNum, service.latestTime, &service.videos)
	if err != nil {
		return err
	}
	//调用utils层，如果用户为登录状态，则更新该视频是否被该用户点赞的状态
	latestTime, _ := utils.FillVideoListFields(service.userId, &service.videos) //不是致命错误，不返回

	//准备好时间戳
	if latestTime != nil {
		//如果latestTime不为nil，说明有视频被点赞了，需要更新nextTime
		service.nextTime = (*latestTime).UnixNano() / 1e6
		return nil
	}
	service.nextTime = time.Now().Unix() / 1e6
	return nil
}

// 3.packData：打包数据
func (service *QueryFeedVideoListService) packData() error {
	service.feedVideo = &FeedVideoList{
		Videos:   service.videos,
		NextTime: service.nextTime,
	}
	return nil
}
