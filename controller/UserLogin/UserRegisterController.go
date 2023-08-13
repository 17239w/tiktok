package userlogin

import (
	"fmt"
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
	rawValue, _ := c.Get("password") // 原始密码,密码通常不会作为查询参数直接附加在URL上进行传递
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
	//打印username和password
	fmt.Println(username, password)
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

	// 返回StatusCodeResponse和UserRegisterResponse
	c.JSON(http.StatusOK, UserRegisterResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 0,
		},
		UserLoginResponse: registerResponse,
	})

}
