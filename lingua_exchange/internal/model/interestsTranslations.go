package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type InterestsTranslations struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	TranslationID  int64  `gorm:"column:translation_id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"translationID"`
	TagID          int64  `gorm:"column:tag_id;type:bigint(20)" json:"tagID"`
	LanguageCode   string `gorm:"column:language_code;type:varchar(5)" json:"languageCode"`
	TranslatedName string `gorm:"column:translated_name;type:varchar(50)" json:"translatedName"`
}
