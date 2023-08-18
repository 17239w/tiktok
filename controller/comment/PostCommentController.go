package comment

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tiktok/models"
	"tiktok/service/comment"

	"github.com/gin-gonic/gin"
)

// PostCommentResponse：增加评论的响应
type PostCommentResponse struct {
	models.StatusCodeResponse //models层的StatusCodeResponse
	*comment.Response         //service层的Response
}

// ProxyPostCommentController：增加评论的代理控制器
type ProxyPostCommentController struct {
	*gin.Context

	videoId     int64
	userId      int64
	commentId   int64
	actionType  int64 //CREATE or DELETE
	commentText string
}

// PostCommentController：增加评论的控制器
func PostCommentController(c *gin.Context) {
	NewProxyPostCommentController(c).Do()
}

// NewProxyPostCommentController：创建ProxyPostCommentController
func NewProxyPostCommentController(c *gin.Context) *ProxyPostCommentController {
	return &ProxyPostCommentController{Context: c}
}

// Do：执行增加评论的操作
func (proxy *ProxyPostCommentController) Do() {
	//1.解析参数
	if err := proxy.parseNum(); err != nil {
		proxy.SendError(err.Error())
		return
	}
	//2.调用service层，增加评论
	commentRes, err := comment.PostComment(proxy.userId, proxy.videoId, proxy.commentId, proxy.actionType, proxy.commentText)
	if err != nil {
		proxy.SendError(err.Error())
		return
	}
	//3.成功返回
	proxy.SendOk(commentRes)
}

// 1.parseNum：解析参数
func (proxy *ProxyPostCommentController) parseNum() error {
	//解析user_Id
	rawUserId, _ := proxy.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	proxy.userId = userId

	//解析video_Id
	rawVideoId := proxy.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	if err != nil {
		return err
	}
	proxy.videoId = videoId

	//根据action_Type解析对应的可选参数
	rawActionType := proxy.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	switch actionType {
	case comment.CREATE:
		//解析comment_text
		proxy.commentText = proxy.Query("comment_text")
	case comment.DELETE:
		//解析comment_id
		proxy.commentId, err = strconv.ParseInt(proxy.Query("comment_id"), 10, 64)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("未定义的行为%d", actionType)
	}
	proxy.actionType = actionType
	return nil
}

// SendError：发送错误
func (proxy *ProxyPostCommentController) SendError(msg string) {
	proxy.JSON(http.StatusOK, PostCommentResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		Response: &comment.Response{},
	})
}

// SendOk：发送成功
func (proxy *ProxyPostCommentController) SendOk(comment *comment.Response) {
	proxy.JSON(http.StatusOK, PostCommentResponse{
		StatusCodeResponse: models.StatusCodeResponse{
			StatusCode: 0,
		},
		Response: comment,
	})
}
