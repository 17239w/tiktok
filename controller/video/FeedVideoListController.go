package video

import (
	"errors"
	"net/http"
	"strconv"
	"tiktok/middleware"
	"tiktok/models"
	"tiktok/service/video"
	"time"

	"github.com/gin-gonic/gin"
)

// FeedResponse：流响应
type FeedResponse struct {
	models.StatusCodeResponse //models层的通用响应
	*video.FeedVideoList      //service层的流视频列表
}

// ProxyFeedVideoList：流视频列表代理
type ProxyFeedVideoList struct {
	*gin.Context
}

// NewProxyFeedVideoList：创建ProxyFeedVideoList
func NewProxyFeedVideoList(c *gin.Context) *ProxyFeedVideoList {
	return &ProxyFeedVideoList{Context: c}
}

// FeedVideoListController：流视频列表控制器
func FeedVideoListController(c *gin.Context) {
	//1.创建ProxyFeedVideoList
	proxy := NewProxyFeedVideoList(c)
	//2.获取token
	token, ok := c.GetQuery("token")
	//无登录状态
	if !ok {
		err := proxy.DoNoToken()
		if err != nil {
			proxy.FeedVideoListError(err.Error())
		}
		return
	}
	//有登录状态
	err := proxy.DoHasToken(token)
	if err != nil {
		proxy.FeedVideoListError(err.Error())
	}
}

// DoNoToken：未登录的视频流推送处理
func (proxy *ProxyFeedVideoList) DoNoToken() error {
	//获取latest_time
	rawTimestamp := proxy.Query("latest_time")
	//获取最新时间
	var latestTime time.Time
	//将rawTimestamp转换为int64
	intTime, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err == nil {
		latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
	}
	//调用service层接口，查询流视频列表(user_id被设置为 0)
	videoList, err := video.QueryFeedVideoList(0, latestTime)
	if err != nil {
		return err
	}
	//查询成功
	proxy.FeedVideoListSuccess(videoList)
	return nil
}

// DoHasToken：如果是登录状态，则生成UserId字段
func (proxy *ProxyFeedVideoList) DoHasToken(token string) error {
	//解析token成功
	if claim, ok := middleware.ParseToken(token); ok {
		//如果token超时
		if time.Now().Unix() > claim.ExpiresAt {
			return errors.New("token超时")
		}
		//获取latest_time
		rawTimestamp := proxy.Query("latest_time")
		//获取最新时间
		var latestTime time.Time
		//将rawTimestamp转换为int64
		intTime, err := strconv.ParseInt(rawTimestamp, 10, 64)
		if err != nil {
			latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
		}
		//调用service层接口,查询流视频列表
		videoList, err := video.QueryFeedVideoList(claim.UserId, latestTime)
		if err != nil {
			return err
		}
		//查询成功
		proxy.FeedVideoListSuccess(videoList)
		return nil
	}
	//解析失败
	return errors.New("token不正确")
}

// FeedVideoListError：流视频列表错误响应
func (proxy *ProxyFeedVideoList) FeedVideoListError(msg string) {
	proxy.JSON(http.StatusOK, FeedResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		}})
}

// FeedVideoListSuccess：流视频列表正确响应
func (proxy *ProxyFeedVideoList) FeedVideoListSuccess(videoList *video.FeedVideoList) {
	proxy.JSON(http.StatusOK, FeedResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 0,
		},
		FeedVideoList: videoList,
	},
	)
}
