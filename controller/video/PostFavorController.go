package video

import (
	"errors"
	"net/http"
	"strconv"
	"tiktok/models"
	"tiktok/service/video"

	"github.com/gin-gonic/gin"
)

// ProxyPostFavorController：代理点赞操作的控制器
type ProxyPostFavorController struct {
	*gin.Context

	userId     int64
	videoId    int64
	actionType int64
}

// PostFavorController：点赞/取消点赞的控制器
func PostFavorController(c *gin.Context) {
	NewProxyPostFavorController(c).Do()
}

// NewProxyPostFavorController：初始化ProxyPostFavorController
func NewProxyPostFavorController(c *gin.Context) *ProxyPostFavorController {
	return &ProxyPostFavorController{Context: c}
}

// Do：执行点赞/取消点赞的handler
func (p *ProxyPostFavorController) Do() {
	//1.解析参数
	if err := p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}
	//2.调用service层，执行点赞/取消点赞的操作
	err := video.PostFavorState(p.userId, p.videoId, p.actionType)
	if err != nil {
		p.SendError(err.Error())
		return
	}
	//3.成功返回
	p.SendOk()
}

// 1.parseNum：解析参数
func (proxy *ProxyPostFavorController) parseNum() error {
	//解析user_Id
	rawUserId, _ := proxy.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	//解析video_Id
	rawVideoId := proxy.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	if err != nil {
		return err
	}
	//解析action_Type
	rawActionType := proxy.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	if err != nil {
		return err
	}
	proxy.videoId = videoId
	proxy.actionType = actionType
	proxy.userId = userId
	return nil
}

// SendError：发送错误信息
func (proxy *ProxyPostFavorController) SendError(msg string) {
	proxy.JSON(http.StatusOK, models.StatusCodeResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

// SendOk：发送成功信息
func (proxy *ProxyPostFavorController) SendOk() {
	proxy.JSON(http.StatusOK, models.StatusCodeResponse{
		StatusCode: 0,
	})
}
