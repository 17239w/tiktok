package userinfo

import (
	"errors"
	"net/http"
	"tiktok/models"

	"github.com/gin-gonic/gin"
)

// UserResponse：用户信息响应
type UserResponse struct {
	models.StatusCodeResponse
	User *models.UserInfo `json:"user"`
}

// 代理模式
type ProxyUserInfo struct {
	context *gin.Context
}

// NewProxyUserInfo：创建代理模式
func NewProxyUserInfo(c *gin.Context) *ProxyUserInfo {
	return &ProxyUserInfo{context: c}
}

// Controller层：负责接收请求，调用models层的方法，返回响应
func UserInfoController(c *gin.Context) {
	proxy := NewProxyUserInfo(c)
	//从JWTMiddleware()中间件中获取user_id
	id, ok := c.Get("user_id")
	if !ok {
		proxy.UserInfoError("获取用户信息失败")
		return
	}
	//查询用户信息
	err := proxy.ControllerQueryUserInfoByUserId(id)
	if err != nil {
		proxy.UserInfoError(err.Error())
		return
	}
}

// ControllerQueryUserInfoByUserId：Controller层，根据user_id查询用户信息
func (proxy *ProxyUserInfo) ControllerQueryUserInfoByUserId(id interface{}) error {
	//将id转换为int64类型
	userId, ok := id.(int64)
	if !ok {
		return errors.New("用户id类型错误")
	}
	// 调用models层方法：根据id查询用户信息
	var userinfo models.UserInfo
	userinfoDAO := models.NewUserInfoDAO()
	err := userinfoDAO.QueryUserInfoById(userId, &userinfo)
	if err != nil {
		return err
	}
	//返回响应
	proxy.UserInfoSuccess(&userinfo)
	return nil
}

// UserInfoError：返回错误信息，查询用户信息失败
func (p *ProxyUserInfo) UserInfoError(msg string) {
	p.context.JSON(http.StatusOK, UserResponse{
		StatusCodeResponse: models.StatusCodeResponse{StatusCode: 1, StatusMsg: msg},
	})
}

// UserInfoSuccess：返回用户信息，查询用户信息成功
func (p *ProxyUserInfo) UserInfoSuccess(user *models.UserInfo) {
	p.context.JSON(http.StatusOK, UserResponse{
		StatusCodeResponse: models.StatusCodeResponse{StatusCode: 0},
		User:               user,
	})
}
