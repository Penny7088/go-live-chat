package model

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// TalkSession 会话列表
type TalkSession struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	TalkType   uint      `gorm:"column:talk_type;type:tinyint(4) unsigned;default:1;NOT NULL" json:"talkType"`     // 聊天类型[1:私信;2:群聊;]
	UserID     int64     `gorm:"column:user_id;type:bigint(20);default:0;NOT NULL" json:"userID"`                  // 用户ID
	ReceiverID uint64    `gorm:"column:receiver_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"receiverID"` // 接收者ID（用户ID 或 群ID）
	IsTop      uint      `gorm:"column:is_top;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isTop"`           // 是否置顶[0:否;1:是;]
	IsDisturb  uint      `gorm:"column:is_disturb;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isDisturb"`   // 消息免打扰[0:否;1:是;]
	IsDelete   uint      `gorm:"column:is_delete;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isDelete"`     // 是否删除[0:否;1:是;]
	IsRobot    uint      `gorm:"column:is_robot;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isRobot"`       // 是否机器人[0:否;1:是;]
	CreatedAt  time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`                                     // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`                                     // 更新时间
}

// TableName table name
func (m *TalkSession) TableName() string {
	return "talk_session"
}

type SearchTalkSession struct {
	Id          int       `json:"id" `
	TalkType    int       `json:"talk_type" `
	ReceiverId  int       `json:"receiver_id" `
	IsDelete    int       `json:"is_delete"`
	IsTop       int       `json:"is_top"`
	IsRobot     int       `json:"is_robot"`
	IsDisturb   int       `json:"is_disturb"`
	UserAvatar  string    `json:"user_avatar"`
	Nickname    string    `json:"nickname"`
	GroupName   string    `json:"group_name"`
	GroupAvatar string    `json:"group_avatar"`
	UpdatedAt   time.Time `json:"updated_at"`
}
