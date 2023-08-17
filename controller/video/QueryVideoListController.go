package video

import (
	"errors"
	"net/http"
	"tiktok/models"
	"tiktok/service/video"

	"github.com/gin-gonic/gin"
)

// ListRespense：视频列表响应
type ListResponse struct {
	models.StatusCodeResponse
	*video.List // service层返回的视频列表
}

// ProxyQueryVideoList：查询视频列表代理
type ProxyQueryVideoList struct {
	context *gin.Context
}

// QueryVideoListController：查询视频列表
func QueryVideoListController(c *gin.Context) {
	// 创建代理
	proxy := NewProxyQueryVideoList(c)
	// 获取userId
	rawId, _ := c.Get("user_id")
	// 根据userId字段进行查询
	err := proxy.DoQueryVideoListByUserId(rawId)
	if err != nil {
		proxy.QueryVideoListError(err.Error())
	}
}

// NewProxyQueryVideoList：创建查询视频列表代理
func NewProxyQueryVideoList(c *gin.Context) *ProxyQueryVideoList {
	return &ProxyQueryVideoList{context: c}
}

// DoQueryVideoListByUserId：根据userId字段进行查询
func (proxy *ProxyQueryVideoList) DoQueryVideoListByUserId(rawId interface{}) error {
	// 解析userId
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	// 调用service层，查询视频列表
	videoList, err := video.QueryVideoListByUserId(userId)
	if err != nil {
		return err
	}
	// 返回响应
	proxy.QueryVideoListOk(videoList)
	return nil
}

// QueryVideoListError：查询视频列表失败
func (proxy *ProxyQueryVideoList) QueryVideoListError(msg string) {
	proxy.context.JSON(http.StatusOK, ListResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		}})
}

// QueryVideoListOk：查询视频列表成功
func (proxy *ProxyQueryVideoList) QueryVideoListOk(videoList *video.List) {
	proxy.context.JSON(http.StatusOK, ListResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 0,
		},
		List: videoList,
	})
}
