package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type Interests struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	TagID   int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"tagID"`
	TagName string `gorm:"column:tag_name;type:varchar(50);NOT NULL" json:"tagName"`
}
