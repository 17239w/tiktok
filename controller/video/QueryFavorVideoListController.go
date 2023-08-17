package video

import (
	"errors"
	"net/http"
	"tiktok/models"
	"tiktok/service/video"

	"github.com/gin-gonic/gin"
)

// FavorVideoListResponse：点赞视频列表的响应
type FavorVideoListResponse struct {
	models.StatusCodeResponse
	*video.FavorList
}

// ProxyFavorVideoListController：点赞列表查询控制器的代理
type ProxyFavorVideoListController struct {
	*gin.Context

	userId int64
}

// QueryFavorVideoListController：点赞列表查询控制器
func QueryFavorVideoListController(c *gin.Context) {
	NewProxyFavorVideoListController(c).Do()
}

// NewProxyFavorVideoListController：ProxyFavorVideoListController构造函数
func NewProxyFavorVideoListController(c *gin.Context) *ProxyFavorVideoListController {
	return &ProxyFavorVideoListController{Context: c}
}

// Do：ProxyFavorVideoListHandler执行函数
func (proxy *ProxyFavorVideoListController) Do() {
	//1.解析参数
	if err := proxy.parseNum(); err != nil {
		proxy.SendError(err.Error())
		return
	}
	//2.调用service层，查询用户的点赞视频列表
	favorVideoList, err := video.QueryFavorVideoList(proxy.userId)
	if err != nil {
		proxy.SendError(err.Error())
		return
	}
	//3.成功返回
	proxy.SendOk(favorVideoList)
}

// 1.解析参数
func (proxy *ProxyFavorVideoListController) parseNum() error {
	rawUserId, _ := proxy.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	proxy.userId = userId
	return nil
}

// 返回错误
func (proxy *ProxyFavorVideoListController) SendError(msg string) {
	proxy.JSON(http.StatusOK, FavorVideoListResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		}})
}

// 返回成功
func (proxy *ProxyFavorVideoListController) SendOk(favorList *video.FavorList) {
	proxy.JSON(http.StatusOK, FavorVideoListResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 0,
		},
		FavorList: favorList,
	})
}
