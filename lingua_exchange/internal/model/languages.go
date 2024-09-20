package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type Languages struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	Name       string `gorm:"column:name;type:varchar(100);NOT NULL" json:"name"`
	NativeName string `gorm:"column:native_name;type:varchar(100);NOT NULL" json:"nativeName"`
	IsoCode    string `gorm:"column:iso_code;type:varchar(10);NOT NULL" json:"isoCode"`
}
