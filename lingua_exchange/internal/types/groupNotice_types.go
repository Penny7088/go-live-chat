package types

import "time"

type GroupNoticeListRequest struct {
	GroupID int `json:"group_id" binding:"required"`
}

type GroupNoticeDeleteRequest struct {
	GroupID  int `json:"group_id" binding:"required"`
	NoticeID int `json:"notice_id" binding:"required"`
}

type GroupNoticeEditRequest struct {
	GroupID   int    `json:"group_id" binding:"required"`
	NoticeID  int    `json:"notice_id"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	IsTop     int    `json:"is_top"`
	IsConfirm int    `json:"is_confirm" `
}

type SearchNoticeItem struct {
	Id           int       `json:"id" grom:"column:id"`
	CreatorId    int       `json:"creator_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	IsTop        int       `json:"is_top"`
	IsConfirm    int       `json:"is_confirm"`
	ConfirmUsers string    `json:"confirm_users"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Avatar       string    `json:"avatar"`
	Nickname     string    `json:"nickname"`
}

// NoticeItemReply 表示项目的结构体
type NoticeItemReply struct {
	ID           int32  `json:"id,"`
	Title        string `json:"title,"`
	Content      string `json:"content,"`
	IsTop        int32  `json:"is_top,"`
	IsConfirm    int32  `json:"is_confirm,"`
	ConfirmUsers string `json:"confirm_users,"`
	Avatar       string `json:"avatar,"`
	CreatorID    int32  `json:"creator_id,"`
	CreatedAt    string `json:"created_at,"`
	UpdatedAt    string `json:"updated_at,"`
}

// TalkRecordExtraGroupNotice 发布群公告
type TalkRecordExtraGroupNotice struct {
	OwnerId   int    `json:"owner_id"`   // 操作人ID
	OwnerName string `json:"owner_name"` // 操作人昵称
	Title     string `json:"title"`      // 标题
	Content   string `json:"content"`    // 内容
}
