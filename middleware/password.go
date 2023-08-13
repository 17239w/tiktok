package middleware

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strconv"

	"tiktok/models"

	"github.com/gin-gonic/gin"
)

// SHA1：对字符串进行SHA1哈希
func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	// 返回16进制编码的字符串
	return hex.EncodeToString(o.Sum(nil))
}

// AuthMiddleWare：鉴权中间件，对密码进行SHA1哈希
func AuthMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 从请求中获取 password
		password := context.Query("password")
		// 如果请求中没有 password，从请求体中获取
		if password == "" {
			password = context.PostForm("password")
		}
		// 将password存入 Gin上下文中
		context.Set("password", SHA1(password))
		context.Next()
	}
}

// NoAuthToGetUserId：无鉴权中间件，从请求中获取用户id
func NoAuthToGetUserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求中获取用户id
		rawId := c.Query("user_id")
		//如果请求中没有用户id，从PostForm中获取
		if rawId == "" {
			rawId = c.PostForm("user_id")
		}
		//用户不存在
		if rawId == "" {
			c.JSON(http.StatusOK, models.StatusCodeResponse{StatusCode: 401, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
			return
		}
		//用户存在，将user_id转换为int64类型
		userId, err := strconv.ParseInt(rawId, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, models.StatusCodeResponse{StatusCode: 401, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
		}
		//将用户id存入Gin上下文
		c.Set("user_id", userId)
		c.Next()
	}
}
