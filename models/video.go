package models

import (
	"errors"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

// Video：视频信息表
type Video struct {
	Id            int64       `json:"id,omitempty"`
	UserInfoId    int64       `json:"-"`                         //作者id
	Author        UserInfo    `json:"author,omitempty" gorm:"-"` //作者v视频=1vn，故gorm不能存他，但json需要返回它
	PlayUrl       string      `json:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url,omitempty"` //封面地址
	FavoriteCount int64       `json:"favorite_count,omitempty"`
	CommentCount  int64       `json:"comment_count,omitempty"`
	IsFavorite    bool        `json:"is_favorite,omitempty"` //是否已经被当前用户点赞
	Title         string      `json:"title,omitempty"`
	Users         []*UserInfo `json:"-" gorm:"many2many:user_favor_videos;"` // 喜欢该视频的用户
	Comments      []*Comment  `json:"-"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
}

// VideoDAO：视频信息表的数据访问对象
type VideoDAO struct {
}

var (
	videoDAO *VideoDAO
	//单例模式
	videoOnce sync.Once
)

// NewVideoDAO：创建VideoDAO
func NewVideoDAO() *VideoDAO {
	videoOnce.Do(func() {
		videoDAO = new(VideoDAO)
	})
	return videoDAO
}

// AddVideo：添加视频
// 注意：由于视频和userinfo有多对一的关系，所以传入的Video参数一定要进行id的映射处理！
func (v *VideoDAO) AddVideo(video *Video) error {
	if video == nil {
		return errors.New("AddVideo 空指针")
	}
	return DB.Create(video).Error
}

// QueryVideoByVideoId：根据视频id查询视频信息
func (v *VideoDAO) QueryVideoByVideoId(videoId int64, video *Video) error {
	if video == nil {
		return errors.New("QueryVideoByVideoId 空指针")
	}
	return DB.Where("id=?", videoId).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title"}).
		First(video).Error
}

// QueryVideoCountByUserId：查询某个用户的视频数量
func (v *VideoDAO) QueryVideoCountByUserId(userId int64, count *int64) error {
	if count == nil {
		return errors.New("QueryVideoCountByUserId 空指针")
	}
	return DB.Model(&Video{}).Where("user_info_id=?", userId).Count(count).Error
}

// QueryVideoListByUserId：查询某个用户的视频列表
func (v *VideoDAO) QueryVideoListByUserId(userId int64, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByUserId videoList 空指针")
	}
	return DB.Where("user_info_id=?", userId).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title"}).
		Find(videoList).Error
}

// GetVideoListByRecentUpload ：返回按投稿时间倒序的视频列表，并限制为最多limit(30)个
func (v *VideoDAO) GetVideoListByRecentUpload(limit int, latestTime time.Time, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByLimit 空指针")
	}
	return DB.Model(&Video{}).Where("created_at<?", latestTime).
		Order("created_at ASC").Limit(limit).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title", "created_at", "updated_at"}).
		Find(videoList).Error
}

// PlusOneFavorByUserIdAndVideoId：增加一个赞
func (v *VideoDAO) PlusOneFavorByUserIdAndVideoId(userId int64, videoId int64) error {
	//数据库事务：更新videos的favorite_count，在user_favor_videos中插入一条记录
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE videos SET favorite_count=favorite_count+1 WHERE id = ?", videoId).Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO `user_favor_videos` (`user_info_id`,`video_id`) VALUES (?,?)", userId, videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

// MinusOneFavorByUserIdAndVideoId：减少一个赞
func (v *VideoDAO) MinusOneFavorByUserIdAndVideoId(userId int64, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		//执行-1之前需要先判断是否合法(不能被减少为负数)
		///数据库事务：更新videos的favorite_count，在user_favor_videos中删除一条记录
		if err := tx.Exec("UPDATE videos SET favorite_count=favorite_count-1 WHERE id = ? AND favorite_count>0", videoId).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM `user_favor_videos`  WHERE `user_info_id` = ? AND `video_id` = ?", userId, videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

// QueryFavorVideoListByUserId：查询某个用户的点赞视频列表
func (v *VideoDAO) QueryFavorVideoListByUserId(userId int64, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryFavorVideoListByUserId videoList 空指针")
	}
	//多表查询，左连接得到结果，再映射到数据
	if err := DB.Raw("SELECT v.* FROM user_favor_videos u , videos v WHERE u.user_info_id = ? AND u.video_id = v.id", userId).Scan(videoList).Error; err != nil {
		return err
	}
	//如果id为0，则说明没有查到数据
	if len(*videoList) == 0 || (*videoList)[0].Id == 0 {
		return errors.New("点赞列表为空")
	}
	return nil
}

// IsVideoExistById：根据视频id判断视频是否存在
func (v *VideoDAO) IsVideoExistById(id int64) bool {
	var video Video
	if err := DB.Where("id=?", id).Select("id").First(&video).Error; err != nil {
		log.Println("models层：IsVideoExistById", err)
	}
	if video.Id == 0 {
		return false
	}
	return true
}
