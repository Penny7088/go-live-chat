package model

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// GroupNotice 群组公告表
type GroupNotice struct {
	ggorm.Model  `gorm:"embedded"` // embed id and time
	ID           uint64            `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	GroupID      uint              `gorm:"column:group_id;type:int(11) unsigned;default:0;NOT NULL" json:"groupID"`        // 群组ID
	CreatorID    uint              `gorm:"column:creator_id;type:int(11) unsigned;default:0;NOT NULL" json:"creatorID"`    // 创建者用户ID
	Title        string            `gorm:"column:title;type:varchar(64);NOT NULL" json:"title"`                            // 公告标题
	Content      string            `gorm:"column:content;type:text;NOT NULL" json:"content"`                               // 公告内容
	ConfirmUsers string            `gorm:"column:confirm_users" json:"confirm_users"`                                      // 已确认成员
	IsDelete     uint              `gorm:"column:is_delete;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isDelete"`   // 是否删除[0:否;1:是;]
	IsTop        uint              `gorm:"column:is_top;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isTop"`         // 是否置顶[0:否;1:是;]
	IsConfirm    uint              `gorm:"column:is_confirm;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isConfirm"` // 是否需群成员确认公告[0:否;1:是;]

	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"` // 更新时间
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`          // 删除时间
}

// TableName table name
func (m *GroupNotice) TableName() string {
	return "group_notice"
}
