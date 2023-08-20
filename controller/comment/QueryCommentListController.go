package comment

import (
	"errors"
	"net/http"
	"strconv"
	"tiktok/controller/video"
	"tiktok/models"
	"tiktok/service/comment"

	"github.com/gin-gonic/gin"
)

// ListResponse：评论列表的响应
type ListResponse struct {
	models.StatusCodeResponse //models.StatusCodeResponse：状态码响应
	*comment.List             //service层返回的评论列表
}

// ProxyCommentListController：评论列表的代理的控制器
type ProxyCommentListController struct {
	*gin.Context

	videoId int64
	userId  int64
}

// QueryCommentListController：查询评论列表的控制器
func QueryCommentListController(c *gin.Context) {
	NewProxyCommentListController(c).Do()
}

// NewProxyCommentListController：创建ProxyCommentListController
func NewProxyCommentListController(c *gin.Context) *ProxyCommentListController {
	return &ProxyCommentListController{Context: c}
}

// Do：执行评论列表的查询
func (proxy *ProxyCommentListController) Do() {
	//1.解析参数
	if err := proxy.parseNum(); err != nil {
		proxy.SendError(err.Error())
		return
	}
	//2.调用service层，查询评论列表
	commentList, err := comment.QueryCommentList(proxy.userId, proxy.videoId)
	if err != nil {
		proxy.SendError(err.Error())
		return
	}
	//3.成功返回
	proxy.SendOk(commentList)
}

// 1.parseNum：解析参数
func (proxy *ProxyCommentListController) parseNum() error {
	//解析user_id
	rawUserId, _ := proxy.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	proxy.userId = userId

	//解析video_id
	rawVideoId := proxy.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	if err != nil {
		return err
	}
	proxy.videoId = videoId

	return nil
}

// SendError：发送错误
func (proxy *ProxyCommentListController) SendError(msg string) {
	proxy.JSON(http.StatusOK, video.FavorVideoListResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		}})
}

// SendOk：发送成功
func (proxy *ProxyCommentListController) SendOk(commentList *comment.List) {
	proxy.JSON(http.StatusOK, ListResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 0,
		},
		List: commentList,
	})
}
