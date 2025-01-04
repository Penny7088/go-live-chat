package model

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// TalkRecordsRead 用户已读列表
type TalkRecordsRead struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	MsgID      string    `gorm:"column:msg_id;type:varchar(64);NOT NULL" json:"msgID"`                          // 消息ID
	UserID     uint      `gorm:"column:user_id;type:int(11) unsigned;default:0;NOT NULL" json:"userID"`         // 用户ID
	ReceiverID uint      `gorm:"column:receiver_id;type:int(11) unsigned;default:0;NOT NULL" json:"receiverID"` // 接受者ID
	CreatedAt  time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`                                  // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`                                  // 更新时间
}

// TableName table name
func (m *TalkRecordsRead) TableName() string {
	return "talk_records_read"
}
