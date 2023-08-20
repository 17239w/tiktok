package userinfo

import (
	"errors"
	"net/http"
	"tiktok/models"
	userinfo "tiktok/service/UserInfo"

	"github.com/gin-gonic/gin"
)

// FollowerListResponse：粉丝列表的响应
type FollowerListResponse struct {
	models.StatusCodeResponse
	*userinfo.FollowerList //service层返回的粉丝列表
}

// ProxyQueryFollowerController：代理查询粉丝列表的控制器
type ProxyQueryFollowerController struct {
	*gin.Context

	userId int64

	*userinfo.FollowerList //service层返回的粉丝列表
}

// QueryFollowerController：查询粉丝列表的控制器
func QueryFollowerController(c *gin.Context) {
	NewProxyQueryFollowerController(c).Do()
}

// NewProxyQueryFollowerController：创建ProxyQueryFollowerController
func NewProxyQueryFollowerController(c *gin.Context) *ProxyQueryFollowerController {
	return &ProxyQueryFollowerController{Context: c}
}

// Do：执行查询粉丝列表
func (proxy *ProxyQueryFollowerController) Do() {
	var err error
	//1.解析参数
	if err = proxy.parseNum(); err != nil {
		proxy.SendError(err.Error())
		return
	}
	//2.准备数据
	if err = proxy.prepareData(); err != nil {
		if errors.Is(err, userinfo.ErrUserNotExist) {
			proxy.SendError(err.Error())
		} else {
			proxy.SendError("准备数据出错")
		}
		return
	}
	proxy.SendOk("成功")
}

// 1.parseNum：解析参数
func (proxy *ProxyQueryFollowerController) parseNum() error {
	rawUserId, _ := proxy.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	proxy.userId = userId
	return nil
}

// 2.prepareData：准备数据
func (proxy *ProxyQueryFollowerController) prepareData() error {
	//调用service层，查询粉丝列表
	list, err := userinfo.QueryFollowerList(proxy.userId)
	if err != nil {
		return err
	}
	proxy.FollowerList = list
	return nil
}

// SendError：发送错误响应
func (proxy *ProxyQueryFollowerController) SendError(msg string) {
	proxy.JSON(http.StatusOK, FollowerListResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

// SendOk：发送成功响应
func (proxy *ProxyQueryFollowerController) SendOk(msg string) {
	proxy.JSON(http.StatusOK, FollowerListResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 0,
			StatusMsg:  msg,
		},
		FollowerList: proxy.FollowerList,
	})
}
