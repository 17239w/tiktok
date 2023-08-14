package router

import (
	"tiktok/config"
	userinfo "tiktok/controller/UserInfo"
	userlogin "tiktok/controller/UserLogin"
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
	// 基础接口
	// /douyin/user/register/：新用户注册时提供用户名，密码即可，用户名需要保证唯一。创建成功后返回用户id和权限token
	baseGroup.POST("/user/register/", middleware.AuthMiddleWare(), userlogin.UserRegisterController)
	// /douyin/user/login/：通过用户名和密码进行登录，登录成功后返回用户id和权限token
	baseGroup.POST("/user/login/", middleware.AuthMiddleWare(), userlogin.UserLoginController)
	// /douyin/user/：获取用户的id、昵称，如果实现社交部分的功能，还会返回关注数和粉丝数
	baseGroup.GET("/user", middleware.JWTMiddleware(), userinfo.UserInfoController)
	// /douyin/publish/action/：登录用户选择视频上传
	//baseGroup.POST("/publish/action/", middleware.JWTMiddleware(), video.PublishVideoController)
	// /douyin/publish/list/：用户的视频发布列表，直接列出用户所有投稿过的视频
	// /douyin/feed/：不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
	// 互动接口
	// 社交接口
	return r
}
