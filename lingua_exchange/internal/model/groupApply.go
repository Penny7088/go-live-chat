package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// GroupApply 群聊成员
type GroupApply struct {
	ggorm.Model `gorm:"embedded"` // embed id and time
	ID          uint64            `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	GroupID     uint              `gorm:"column:group_id;type:int(11) unsigned;default:0;NOT NULL" json:"groupID"` // 群组ID
	UserID      uint              `gorm:"column:user_id;type:int(11) unsigned;default:0;NOT NULL" json:"userID"`   // 用户ID
	Status      int               `gorm:"column:status;type:int(11);default:1;NOT NULL" json:"status"`             // 申请状态
	Remark      string            `gorm:"column:remark;type:varchar(255);NOT NULL" json:"remark"`                  // 备注信息
	Reason      string            `gorm:"column:reason;type:varchar(255);NOT NULL" json:"reason"`                  // 拒绝原因
}

// TableName table name
func (m *GroupApply) TableName() string {
	return "group_apply"
}
