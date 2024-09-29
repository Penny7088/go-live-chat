package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
	"time"
)

type UserDevices struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	UserID      int64     `gorm:"column:user_id;type:bigint(20);NOT NULL" json:"userId"`
	DeviceToken string    `gorm:"column:device_token;type:varchar(255);NOT NULL" json:"deviceToken"`
	DeviceType  string    `gorm:"column:device_type;type:enum('iOS','Android','Web');NOT NULL" json:"deviceType"`
	IPAddress   string    `gorm:"column:ip_address;type:varchar(255);NOT NULL" json:"ipAddress"`
	LastActive  time.Time `gorm:"column:last_active;type:timestamp;default:CURRENT_TIMESTAMP" json:"lastActive"`
}
