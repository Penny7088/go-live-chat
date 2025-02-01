package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// Group 用户聊天群
type Group struct {
	ggorm.Model `gorm:"embedded"` // embed id and time
	ID          int               `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`                                 // 群ID
	Type        uint              `gorm:"column:type;type:tinyint(4) unsigned;default:1;NOT NULL" json:"type"`            // 群类型[1:普通群;2:企业群;]
	Name        string            `gorm:"column:name;type:varchar(64);NOT NULL" json:"name"`                              // 群名称
	Profile     string            `gorm:"column:profile;type:varchar(128);NOT NULL" json:"profile"`                       // 群介绍
	Avatar      string            `gorm:"column:avatar;type:varchar(255);NOT NULL" json:"avatar"`                         // 群头像
	MaxNum      uint              `gorm:"column:max_num;type:smallint(6) unsigned;default:200;NOT NULL" json:"maxNum"`    // 最大群成员数量
	IsOvert     uint              `gorm:"column:is_overt;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isOvert"`     // 是否公开可见[0:否;1:是;]
	IsMute      uint              `gorm:"column:is_mute;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isMute"`       // 是否全员禁言 [0:否;1:是;]，提示:不包含群主或管理员
	IsDismiss   uint              `gorm:"column:is_dismiss;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isDismiss"` // 是否已解散[0:否;1:是;]
	CreatorID   int               `gorm:"column:creator_id;type:int(11) unsigned;default:0;NOT NULL" json:"creatorID"`    // 创建者ID(群主ID)
}

// TableName table name
func (m *Group) TableName() string {
	return "group"
}
