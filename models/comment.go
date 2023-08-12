package models

import (
	"time"
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
