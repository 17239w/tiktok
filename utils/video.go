package utils

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"tiktok/cache"
	"tiktok/config"
	"tiktok/models"
	"time"
)

// GetFileUrl：获取文件的url
func GetFileUrl(fileName string) string {
	base := fmt.Sprintf("http://%s:%d/static/%s", config.Global.IP, config.Global.Port, fileName)
	return base
}

// NewFileName：根据user_id+用户发布的视频数量连接成独一无二的文件名
func NewFileName(userId int64) string {
	var count int64
	//调用models层的方法,获取用户发布的视频数量
	err := models.NewVideoDAO().QueryVideoCountByUserId(userId, &count)
	if err != nil {
		log.Println(err)
	}
	//返回的文件名格式为：userId-视频数量
	return fmt.Sprintf("%d-%d", userId, count)
}

// FillVideoListFields：填充每个视频的作者信息(因为作者与视频的一对多关系，数据库中存下的是user_id)
// 当user_id>0时，我们判断当前为登录状态，其余情况为未登录状态(不需要填充IsFavorite字段)
func FillVideoListFields(userId int64, videos *[]*models.Video) (*time.Time, error) {
	//判断videos是否为空
	size := len(*videos)
	if videos == nil || size == 0 {
		return nil, errors.New("utils层：FillVideoListFields videos为空")
	}
	//1.获取最近的投稿时间
	latestTime := (*videos)[size-1].CreatedAt
	//2.添加作者信息
	for i := 0; i < size; i++ {
		var userInfo models.UserInfo
		//调用models层，根据用户id获取用户信息
		dao := models.NewUserInfoDAO()
		err := dao.QueryUserInfoById((*videos)[i].UserInfoId, &userInfo)
		if err != nil {
			continue
		}
		//调用cache层，获取用户的粉丝数
		proxy := cache.NewProxyIndexMap()
		userInfo.IsFollow = proxy.GetUserRelation(userId, userInfo.Id) //根据cache更新是否被点赞
		(*videos)[i].Author = userInfo
		//若用户登录，则填充点赞状态(IsFavorite字段)
		if userId > 0 {
			//调用models层，获取视频的点赞状态
			(*videos)[i].IsFavorite = proxy.GetVideoFavorState(userId, (*videos)[i].Id)
		}
	}
	return &latestTime, nil
}

// SaveImageFromVideo：将视频切一帧保存到本地
// isDebug：控制是否打印出执行的ffmepg命令
func SaveImageFromVideo(name string, isDebug bool) error {
	//创建一个Video2Image对象
	changevideoToimage := NewChangeVideoToImage()
	//打印出执行的ffmepg命令
	if isDebug {
		//调用Debug方法
		changevideoToimage.Debug()
	}
	//输入路径
	changevideoToimage.InputPath = filepath.Join(config.Global.StaticSourcePath, name+defaultVideoSuffix)
	//输出路径
	changevideoToimage.OutputPath = filepath.Join(config.Global.StaticSourcePath, name+defaultImageSuffix)
	//截取1帧
	changevideoToimage.FrameCount = 1
	//获取ffmpeg命令行
	queryString, err := changevideoToimage.GetQueryString()
	if err != nil {
		return err
	}
	//执行ffmpeg命令行
	return changevideoToimage.ExecCommand(queryString)
}
