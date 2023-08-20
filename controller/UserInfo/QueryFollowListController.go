package userinfo

import (
	"errors"
	"net/http"
	"tiktok/models"
	userinfo "tiktok/service/UserInfo"

	"github.com/gin-gonic/gin"
)

// FollowListResponse：关注列表的响应
type FollowListResponse struct {
	models.StatusCodeResponse
	*userinfo.FollowList
}

// ProxyQueryFollowList：查询关注列表的代理
type ProxyQueryFollowList struct {
	*gin.Context

	userId int64

	*userinfo.FollowList //关注列表
}

// QueryFollowListController：查询关注列表的控制器
func QueryFollowListController(c *gin.Context) {
	NewProxyQueryFollowList(c).Do()
}

// NewProxyQueryFollowList：创建ProxyQueryFollowList
func NewProxyQueryFollowList(c *gin.Context) *ProxyQueryFollowList {
	return &ProxyQueryFollowList{Context: c}
}

// Do：执行查询关注列表的代理
func (proxy *ProxyQueryFollowList) Do() {
	var err error
	//1.解析参数
	if err = proxy.parseNum(); err != nil {
		proxy.SendError(err.Error())
		return
	}
	//2.准备数据
	if err = proxy.prepareData(); err != nil {
		proxy.SendError(err.Error())
		return
	}
	proxy.SendOk("请求成功")
}

// 1.parseNum：解析参数
func (proxy *ProxyQueryFollowList) parseNum() error {
	rawUserId, _ := proxy.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	proxy.userId = userId
	return nil
}

// 2.prepareData：准备数据
func (proxy *ProxyQueryFollowList) prepareData() error {
	//调用service层，查询用户关注列表
	list, err := userinfo.QueryFollowList(proxy.userId)
	if err != nil {
		return err
	}
	proxy.FollowList = list
	return nil
}

// SendError：发送错误
func (proxy *ProxyQueryFollowList) SendError(msg string) {
	proxy.JSON(http.StatusOK, FollowListResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

// SendOk：发送成功
func (proxy *ProxyQueryFollowList) SendOk(msg string) {
	proxy.JSON(http.StatusOK, FollowListResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 0,
			StatusMsg:  msg,
		},
		FollowList: proxy.FollowList,
	})
}
