package model

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// GroupMember 群聊成员
type GroupMember struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	GroupID  uint      `gorm:"column:group_id;type:int(11) unsigned;default:0;NOT NULL" json:"groupID"`  // 群组ID
	UserID   int       `gorm:"column:user_id;type:int(11) unsigned;default:0;NOT NULL" json:"userID"`    // 用户ID
	Leader   int       `gorm:"column:leader;type:tinyint(4) unsigned;default:0;NOT NULL" json:"leader"`  // 成员属性[0:普通成员;1:管理员;2:群主;]
	UserCard string    `gorm:"column:user_card;type:varchar(64);NOT NULL" json:"userCard"`               // 群名片
	IsQuit   int       `gorm:"column:is_quit;type:tinyint(4);default:0;NOT NULL" json:"isQuit"`          // 是否退群[0:否;1:是;]
	IsMute   uint      `gorm:"column:is_mute;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isMute"` // 是否禁言[0:否;1:是;]
	JoinTime time.Time `gorm:"column:join_time;type:datetime" json:"joinTime"`                           // 入群时间

	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"` // 更新时间
}

// TableName table name
func (m *GroupMember) TableName() string {
	return "group_member"
}

type CountGroupMember struct {
	GroupId int `gorm:"column:group_id;"`
	Count   int `gorm:"column:count;"`
}
