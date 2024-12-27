package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type UserLanguages struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	UserID        int64  `gorm:"column:user_id;type:bigint(20);primary_key" json:"userID"`
	LanguageID    int64  `gorm:"column:language_id;type:bigint(20);NOT NULL" json:"languageID"`
	LanguageLevel string `gorm:"column:language_level;type:varchar(50)" json:"languageLevel"`
}
