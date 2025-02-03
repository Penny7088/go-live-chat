package types

import "time"

// GroupApplyCreateRequest 入群申请
type GroupApplyCreateRequest struct {
	GroupID int    `json:"group_id" binding:"required"`
	Remark  string `json:"remark" binding:"required"`
}

type GroupApplyDeleteRequest struct {
	ApplyID int `json:"apply_id" binding:"required"`
}

type GroupApplyAgreeRequest struct {
	ApplyID int `json:"apply_id" binding:"required"`
}

type GroupApplyDeclineRequest struct {
	ApplyID int    `json:"apply_id" binding:"required"`
	Remark  string `json:"remark" binding:"required"`
}

type GroupApplyListRequest struct {
	GroupID int `json:"group_id" binding:"required"`
}

// GroupApplyListResponse 表示群组申请列表响应
type GroupApplyListResponse struct {
	Items []GroupApplyItem `json:"items,omitempty"`
}

// GroupApplyItem 表示群组申请项
type GroupApplyItem struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
	Remark    string `json:"remark"`
	Avatar    string `json:"avatar"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type GroupApplyList struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`         // 自增ID
	GroupId   int       `gorm:"column:group_id;default:0;NOT NULL" json:"group_id"`     // 群组ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`       // 用户ID
	Remark    string    `gorm:"column:remark;NOT NULL" json:"remark"`                   // 备注信息
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`           // 创建时间
	Nickname  string    `gorm:"column:username;NOT NULL" json:"username"`               // 用户昵称
	Avatar    string    `gorm:"column:profile_picture;NOT NULL" json:"profile_picture"` // 用户头像地址
}

// GroupApplyAllResponse 表示群组申请列表响应
type GroupApplyAllResponse struct {
	Items []GroupApplyItem `json:"items,omitempty"`
}
