package model

type WSServer struct {
	ServerName string `gorm:"column:name;type:varchar(100);NOT NULL" json:"serverName"`
}
