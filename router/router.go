package router

import (
	"tiktok/config"
	userinfo "tiktok/controller/UserInfo"
	"tiktok/middleware"
	"tiktok/models"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	//初始化配置
	models.InitDB()
	r := gin.Default()
	//设置静态资源目录
	r.Static("/static", config.Global.StaticSourcePath)
	//路由组
	baseGroup := r.Group("/douyin")
	// basic apis
	// /douyin/user/：获取用户的id、昵称，如果实现社交部分的功能，还会返回关注数和粉丝数
	baseGroup.GET("/user", middleware.AuthMiddleWare(), userinfo.UserInfoController)
	// extend 1
	// extend 2
	return r
}
