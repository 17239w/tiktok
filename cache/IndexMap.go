package cache

import (
	"context"
	"fmt"
	"tiktok/config"

	"github.com/go-redis/redis/v8"
)

var redisCtx = context.Background() //redis的上下文
var redisClient *redis.Client       //redis的客户端

const (
	favor    = "favor"
	relation = "relation"
)

// init：初始化redis客户端
func init() {
	redisClient = redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Global.RDB.IP, config.Global.RDB.Port),
			Password: "", //没有设置密码
			DB:       config.Global.RDB.Database,
		})
}

// ProxyIndexMap：代理对象
type ProxyIndexMap struct {
}

var (
	proxyIndexOperation ProxyIndexMap
)

// NewProxyIndexMap：创建代理对象proxyIndexOperation
func NewProxyIndexMap() *ProxyIndexMap {
	return &proxyIndexOperation
}

// UpdateVideoFavorState：更新视频点赞状态(state:true为点赞，false为取消点赞)
func (i *ProxyIndexMap) UpdateVideoFavorState(userId int64, videoId int64, state bool) {
	// 存储用户点赞过的video集合,key的格式为：favor:userId
	key := fmt.Sprintf("%s:%d", favor, userId)
	//如果state为true，那么就将videoId添加到redis集合中，否则就将videoId从redis集合中删除
	if state {
		//redisCtx为redis的上下文，key为redis的键，videoId为redis的值
		redisClient.SAdd(redisCtx, key, videoId)
		return
	}
	redisClient.SRem(redisCtx, key, videoId)
}

// GetVideoFavorState：得到点赞状态
func (i *ProxyIndexMap) GetVideoFavorState(userId int64, videoId int64) bool {
	key := fmt.Sprintf("%s:%d", favor, userId)
	//判断videoId是否在redis集合中
	ret := redisClient.SIsMember(redisCtx, key, videoId)
	//ret.Val()为bool类型
	return ret.Val()
}

// UpdateUserRelation：更新用户的关注状态(state:true为关注，false为取消关注)
func (i *ProxyIndexMap) UpdateUserRelation(userId int64, followId int64, state bool) {
	//key的格式为：relation:userId
	key := fmt.Sprintf("%s:%d", relation, userId)
	//state为true，将followId添加到redis集合中
	if state {
		redisClient.SAdd(redisCtx, key, followId)
		return
	}
	//state为false,将followId从redis集合中删除
	redisClient.SRem(redisCtx, key, followId)
}

// GetUserRelation：得到用户关注的状态
func (i *ProxyIndexMap) GetUserRelation(userId int64, followId int64) bool {
	//key的格式为：relation:userId
	key := fmt.Sprintf("%s:%d", relation, userId)
	//判断followId是否在redis集合中
	ret := redisClient.SIsMember(redisCtx, key, followId)
	return ret.Val()
}
