package userlogin

import (
	"net/http"
	"tiktok/models"
	userlogin "tiktok/service/UserLogin"

	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	models.StatusCodeResponse
	*userlogin.UserLoginResponse
}

// UserLoginController：用户登录控制器
func UserLoginHandler(c *gin.Context) {
	username := c.Query("username")
	raw, _ := c.Get("password")
	password, ok := raw.(string)
	if !ok {
		c.JSON(http.StatusOK, UserLoginResponse{
			StatusCodeResponse: models.StatusCodeResponse{
				StatusCode: 1,
				StatusMsg:  "密码解析错误",
			},
		})
	}
	//调用service层的函数,查询用户是否存在
	userLoginResponse, err := userlogin.QueryUserLogin(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			StatusCodeResponse: models.StatusCodeResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		StatusCodeResponse: models.StatusCodeResponse{StatusCode: 0},
		UserLoginResponse:  userLoginResponse,
	})
}
