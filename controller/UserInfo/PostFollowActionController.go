package userinfo

import (
	"errors"
	"net/http"
	"strconv"
	"tiktok/models"
	userinfo "tiktok/service/UserInfo"

	"github.com/gin-gonic/gin"
)

// ProxyPostFollowAction：代理对象
type ProxyPostFollowAction struct {
	*gin.Context

	userId     int64
	followId   int64
	actionType int
}

// PostFollowActionController：关注/取消关注的控制器
func PostFollowActionController(c *gin.Context) {
	NewProxyPostFollowAction(c).Do()
}

// NewProxyPostFollowAction：创建代理对象
func NewProxyPostFollowAction(c *gin.Context) *ProxyPostFollowAction {
	return &ProxyPostFollowAction{Context: c}
}

// Do：执行 关注/取消关注 的操作
func (proxy *ProxyPostFollowAction) Do() {
	var err error
	//1.解析参数
	if err = proxy.prepareNum(); err != nil {
		proxy.SendError(err.Error())
		return
	}
	//2.执行关注/取消关注的操作
	if err = proxy.startAction(); err != nil {
		//当错误为service层发生的，那么就是 用户不存在ErrIvdFolUsr 或者 未定义操作ErrIvdAct
		if errors.Is(err, userinfo.ErrIvdAct) || errors.Is(err, userinfo.ErrIvdFolUsr) {
			proxy.SendError(err.Error())
		} else {
			//当错误为model层发生的，就是重复键值的插入
			proxy.SendError("请勿重复关注")
		}
		return
	}
	proxy.SendOk("操作成功")
}

// 1.prepareNum：解析参数 user_id、to_user_id、action_type
func (proxy *ProxyPostFollowAction) prepareNum() error {
	//解析user_id
	rawUserId, _ := proxy.Get("user_id")
	//将interface{}转换为int64
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	proxy.userId = userId
	//解析to_user_id
	followId := proxy.Query("to_user_id")
	//将string转换为int64
	parseInt, err := strconv.ParseInt(followId, 10, 64)
	if err != nil {
		return err
	}
	proxy.followId = parseInt
	//解析action_type
	actionType := proxy.Query("action_type")
	//将string转换为int
	parseInt, err = strconv.ParseInt(actionType, 10, 32)
	if err != nil {
		return err
	}
	//将int转换为int64
	proxy.actionType = int(parseInt)
	return nil
}

// 2.startAction：执行关注/取消关注的操作
func (proxy *ProxyPostFollowAction) startAction() error {
	//调用service层，执行关注/取消关注的操作
	err := userinfo.PostFollowAction(proxy.userId, proxy.followId, proxy.actionType)
	if err != nil {
		return err
	}
	return nil
}

// SendError：发送错误信息
func (proxy *ProxyPostFollowAction) SendError(msg string) {
	proxy.JSON(http.StatusOK, models.StatusCodeResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

// SendOk：发送成功信息
func (proxy *ProxyPostFollowAction) SendOk(msg string) {
	proxy.JSON(http.StatusOK, models.StatusCodeResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}
