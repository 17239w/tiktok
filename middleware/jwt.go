package middleware

import (
	"log"
	"net/http"
	"tiktok/models" //引用tiktok/models包中的类型、函数和变量
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 包含JWT签名密钥的字节切片
var jwtSecret = []byte("secret")

// JWTClaims：JWT的声明结构
type JWTClaims struct {
	UserId int64
	jwt.StandardClaims
}

// ReleaseToken：颁发token
func ReleaseToken(user models.UserLogin) (string, error) {
	//	设置token过期时间为7天后
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	//	创建声明
	claims := &JWTClaims{
		UserId: user.UserInfoId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),     //发放时间
			Issuer:    "douyin",              //发放者
			Subject:   "L_B__",               //主题
		}}

	//	创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//	签名字符串
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken：解析token
func ParseToken(tokenString string) (*JWTClaims, bool) {
	//	解析token
	token, _ := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	//	判断token是否有效
	if token != nil {
		//	claims-->JWTClaims类型
		if key, ok := token.Claims.(*JWTClaims); ok {
			if token.Valid {
				return key, true
			} else {
				return key, false
			}
		}
	}
	return nil, false
}

// JWTMiddleware：JWT中间件，鉴权并设置user_id
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		tokenStr := c.Query("token")
		//打印日志
		log.Println("middleware层(鉴权并设置user_id):tokenStr:", tokenStr)
		// 如果查询参数中没有token，则尝试从POST表单中获取token字段的值
		// 从POST请求的urlencoded表单或者multipart表单中获取指定的键值。如果键存在，它会返回对应的值；如果键不存在，它会返回一个空字符串（""）
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, models.StatusCodeResponse{StatusCode: 401, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
			return
		}
		//验证token
		tokenStruck, ok := ParseToken(tokenStr)
		if !ok {
			c.JSON(http.StatusOK, models.StatusCodeResponse{
				StatusCode: 403,
				StatusMsg:  "token不正确",
			})
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, models.StatusCodeResponse{
				StatusCode: 402,
				StatusMsg:  "token过期",
			})
			c.Abort() //阻止执行
			return
		}
		//设置user_id
		c.Set("user_id", tokenStruck.UserId)
		c.Next()
	}
}
