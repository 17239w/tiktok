package models

import (
	"errors"
	"log"
	"sync"

	"gorm.io/gorm"
)

// 定义错误
var (
	ErrIvdPtr        = errors.New("invalid pointer")    //无效指针
	ErrEmptyUserList = errors.New("user list is empty") //空用户列表
)

// UserInfo：用户信息
// gorm:"id,omitempty"：表示在将结构体转换为JSON或进行数据库操作时，如果该字段的值为nil/0，则忽略该字段
// json:"-"：表示该字段不会包含在生成的JSON中，也不会参与JSON解码过程
type UserInfo struct {
	Id             int64       `json:"id" gorm:"id,omitempty"`
	Name           string      `json:"name" gorm:"name,omitempty"`
	FollowingCount int64       `json:"following_count" gorm:"following_count,omitempty"` //关注数
	FollowerCount  int64       `json:"follower_count" gorm:"follower_count,omitempty"`   //粉丝数
	IsFollow       bool        `json:"is_follow" gorm:"is_follow,omitempty"`
	User           *UserLogin  `json:"-"`                                     //用户v账号密码=1v1
	Videos         []*Video    `json:"-"`                                     //用户v投稿视频=1vn
	Follows        []*UserInfo `json:"-" gorm:"many2many:user_relations;"`    //用户v关注用户=nvn
	FavorVideos    []*Video    `json:"-" gorm:"many2many:user_favor_videos;"` //用户v点赞视频=nvn
	Comments       []*Comment  `json:"-"`                                     //用户v评论视频=1vn
}

// UserInfoDAO：用户信息DAO
type UserInfoDAO struct {
}

var (
	userInfoDAO *UserInfoDAO
	//单例模式
	userInfoOnce sync.Once
)

// NewUserInfoDAO：创建UserInfoDAO
func NewUserInfoDAO() *UserInfoDAO {
	userInfoOnce.Do(func() {
		userInfoDAO = new(UserInfoDAO)
	})
	return userInfoDAO
}

// QueryUserInfoById：根据用户id查询用户信息
func (dao *UserInfoDAO) QueryUserInfoById(userId int64, userInfo *UserInfo) error {
	//判断指针是否为空
	if userInfo == nil {
		return ErrIvdPtr
	}
	//查询用户信息
	DB.Where("id=?", userId).Select([]string{"id", "name", "following_count", "follower_count"}).First(userInfo)
	//判断用户是否存在
	if userInfo.Id == 0 {
		return errors.New("user is not exist")
	}
	return nil
}

// AddUserInfo：添加用户信息
func (dao *UserInfoDAO) AddUserInfo(userInfo *UserInfo) error {
	//判断指针是否为空
	if userInfo == nil {
		return ErrIvdPtr
	}
	//添加用户信息
	return DB.Create(userInfo).Error
}

// IsUserExistById：判断用户是否存在
func (dao *UserInfoDAO) IsUserExistById(userId int64) bool {
	var userInfo UserInfo
	//查询用户信息
	err := DB.Where("id=?", userId).Select([]string{"id"}).First(&userInfo).Error
	if err != nil {
		log.Println(err)
	}
	if userInfo.Id == 0 {
		return false
	}
	return true
}

// AddUserFollow：添加用户关注
func (dao *UserInfoDAO) AddUserFollow(userId int64, followToId int64) error {
	//判断用户是否存在
	if !dao.IsUserExistById(userId) || !dao.IsUserExistById(followToId) {
		return errors.New("user is not exist")
	}
	//开启事务
	return DB.Transaction(func(tx *gorm.DB) error {
		//添加用户关注
		if err := tx.Exec("UPDATE user_infos SET following_count=following_count+1 WHERE id=?", userId).Error; err != nil {
			return err
		}
		//添加粉丝数
		if err := tx.Exec("UPDATE user_infos SET follower_count=follower_count+1 WHERE id=?", followToId).Error; err != nil {
			return err
		}
		//更新user_relations表
		if err := tx.Exec("INSERT INTO user_relations(user_info_id,follow_id) VALUES(?,?)", userId, followToId).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteUserFollow：删除用户关注
func (dao *UserInfoDAO) CancelUserFollow(userId int64, followToId int64) error {
	//判断用户是否存在
	if !dao.IsUserExistById(userId) || !dao.IsUserExistById(followToId) {
		return errors.New("user is not exist")
	}
	//开启事务
	return DB.Transaction(func(tx *gorm.DB) error {
		//删除用户关注
		if err := tx.Exec("UPDATE user_infos SET following_count=following_count-1 WHERE id=?", userId).Error; err != nil {
			return err
		}
		//删除粉丝数
		if err := tx.Exec("UPDATE user_infos SET follower_count=follower_count-1 WHERE id=?", followToId).Error; err != nil {
			return err
		}
		//更新user_relations表
		if err := tx.Exec("DELETE FROM user_relations WHERE user_info_id=? AND follow_id=?", userId, followToId).Error; err != nil {
			return err
		}
		return nil
	})
}

// QueryFollowListByUserId：查询用户关注列表
func (dao *UserInfoDAO) QueryFollowListByUserId(userId int64, userFollowsList *[]*UserInfo) error {
	// 判断指针是否为空
	if userFollowsList == nil {
		return ErrIvdPtr
	}
	var err error
	// user_relations.user_info_id=? AND user_relations.follow_id=user_infos.id
	if err = DB.Raw("SELECT i.* FROM user_relations r,user_infos i WHERE r.user_info_id=? AND r.follow_id=i.id", userId).Scan(userFollowsList).Error; err != nil {
		return err
	}
	// 若用户关注列表为空
	if len(*userFollowsList) == 0 || (*userFollowsList)[0].Id == 0 {
		return ErrEmptyUserList
	}
	return nil
}

// QueryFollowerListByUserId：查询用户粉丝列表
func (dao *UserInfoDAO) QueryFollowerListByUserId(userId int64, userFollowersList *[]*UserInfo) error {
	// 判断指针是否为空
	if userFollowersList == nil {
		return ErrIvdPtr
	}
	var err error
	// user_relations.follow_id=? AND user_relations.user_info_id=user_infos.id
	if err = DB.Raw("SELECT i.* FROM user_relations r,user_infos i WHERE r.follow_id=? AND r.user_info_id=i.id", userId).Scan(userFollowersList).Error; err != nil {
		return err
	}
	// 若用户粉丝列表为空
	if len(*userFollowersList) == 0 || (*userFollowersList)[0].Id == 0 {
		return ErrEmptyUserList
	}
	return nil
}
