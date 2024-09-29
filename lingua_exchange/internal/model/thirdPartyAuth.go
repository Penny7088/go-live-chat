package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type ThirdPartyAuth struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	UserID         int64  `gorm:"column:user_id;type:bigint(20);NOT NULL" json:"userID"`
	Provider       string `gorm:"column:provider;type:varchar(50);NOT NULL" json:"provider"`
	ProviderUserID string `gorm:"column:provider_user_id;type:varchar(255);NOT NULL" json:"providerUserID"`
}

// TableName table name
func (m *ThirdPartyAuth) TableName() string {
	return "third_party_auth"
}
