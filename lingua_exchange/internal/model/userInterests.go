package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type UserInterests struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	UserID int64 `gorm:"column:user_id;type:bigint(20);primary_key" json:"userID"`
	TagID  int64 `gorm:"column:tag_id;type:bigint(20);NOT NULL" json:"tagID"`
}
