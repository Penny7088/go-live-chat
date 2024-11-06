package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type Countries struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	Name      string `gorm:"column:name;type:varchar(100);NOT NULL" json:"name"`
	IsoCode   string `gorm:"column:iso_code;type:varchar(10);NOT NULL" json:"isoCode"`
	VisitName string `gorm:"column:visit_name;type:varchar(100)" json:"visitName"` // 方便阅读的字段
	PhoneCode int    `gorm:"column:phone_code;type:int(11)" json:"phoneCode"`      // 国家号前缀
}
