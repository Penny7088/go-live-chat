package model

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// TalkRecordsDelete 聊天记录删除记录表
type TalkRecordsDelete struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	MsgID     string    `gorm:"column:msg_id;type:varchar(64);NOT NULL" json:"msgID"`                  // 聊天记录ID
	UserID    uint      `gorm:"column:user_id;type:int(11) unsigned;default:0;NOT NULL" json:"userID"` // 用户ID
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`                          // 创建时间
}

// TableName table name
func (m *TalkRecordsDelete) TableName() string {
	return "talk_records_delete"
}
