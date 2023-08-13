package userlogin

import (
	"net/http"
	"tiktok/models"
	userlogin "tiktok/service/UserLogin"

	"github.com/gin-gonic/gin"
)

// 用户注册响应体
type UserRegisterResponse struct {
	models.StatusCodeResponse
	*userlogin.UserLoginResponse // 保存 用户登陆成功后的相应数据，包括token和id
}

// UserRegisterController：用户注册控制器
func UserRegisterController(c *gin.Context) {
	username := c.Query("username")
	rawValue, _ := c.Get("password") // 原始密码
	// 将password-->string
	password, ok := rawValue.(string)
	if !ok {
		c.JSON(http.StatusOK, UserRegisterResponse{
			StatusCodeResponse: models.StatusCodeResponse{
				StatusCode: 1,
				StatusMsg:  "密码解析错误",
			},
		})
		return
	}
	// 调用service层，进行注册
	registerResponse, err := userlogin.UserRegister(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			StatusCodeResponse: models.StatusCodeResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	// 返回UserRegisterResponse
	c.JSON(http.StatusOK, UserRegisterResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 0,
			StatusMsg:  "注册成功",
		},
		UserLoginResponse: registerResponse,
	})

}
