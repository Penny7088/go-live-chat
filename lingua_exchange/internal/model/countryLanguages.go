package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type CountryLanguages struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	CountryID  int64 `gorm:"column:country_id;type:bigint(20);NOT NULL" json:"countryID"`
	LanguageID int64 `gorm:"column:language_id;type:bigint(20);NOT NULL" json:"languageID"`
}
