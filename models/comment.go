package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Comment：评论
type Comment struct {
	Id         int64     `json:"id"`
	UserInfoId int64     `json:"-"` //用户v评论=1vn      1个用户可以发表多个评论
	VideoId    int64     `json:"-"` //视频v评论=1vn       1个视频可以有多个评论
	User       UserInfo  `json:"user" gorm:"-"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"-"`
	CreateDate string    `json:"create_date" gorm:"-"`
}

// CommentDAO：评论DAO
type CommentDAO struct {
}

var (
	commentDAO CommentDAO
)

// NewCommentDAO：实例化评论DAO
func NewCommentDAO() *CommentDAO {
	return &commentDAO
}

// AddCommentAndUpdateCount：添加评论并更新视频的评论数
func (c *CommentDAO) AddCommentAndUpdateCount(comment *Comment) error {
	if comment == nil {
		return errors.New("AddCommentAndUpdateCount comment空指针")
	}
	// 执行事务 (tx：事务对象)
	return DB.Transaction(func(tx *gorm.DB) error {
		// 添加评论数据
		if err := tx.Create(comment).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		// comment_count+1
		if err := tx.Exec("UPDATE videos v SET v.comment_count = v.comment_count+1 WHERE v.id=?", comment.VideoId).Error; err != nil {
			return err
		}
		// 返回nil提交事务
		return nil
	})
}

// DeleteCommentAndUpdateCountById：根据评论id删除评论并更新视频的评论数
func (c *CommentDAO) DeleteCommentAndUpdateCountById(commentId, videoId int64) error {
	//执行事务
	return DB.Transaction(func(tx *gorm.DB) error {
		//删除评论
		if err := tx.Exec("DELETE FROM comments WHERE id = ?", commentId).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		//comment_count-1
		if err := tx.Exec("UPDATE videos v SET v.comment_count = v.comment_count-1 WHERE v.id=? AND v.comment_count>0", videoId).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

// QueryCommentById：根据评论id查询评论
func (c *CommentDAO) QueryCommentById(id int64, comment *Comment) error {
	if comment == nil {
		return errors.New("QueryCommentById comment 空指针")
	}
	return DB.Where("id=?", id).First(comment).Error
}

// QueryCommentListByVideoId：根据视频id查询评论列表
func (c *CommentDAO) QueryCommentListByVideoId(videoId int64, comments *[]*Comment) error {
	if comments == nil {
		return errors.New("QueryCommentListByVideoId comments空指针")
	}
	if err := DB.Model(&Comment{}).Where("video_id=?", videoId).Find(comments).Error; err != nil {
		return err
	}
	return nil
}
