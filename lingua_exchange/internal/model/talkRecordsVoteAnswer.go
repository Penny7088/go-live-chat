package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// TalkRecordsVoteAnswer 聊天对话记录（投票消息统计表）
type TalkRecordsVoteAnswer struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	VoteID uint   `gorm:"column:vote_id;type:int(11) unsigned;default:0;NOT NULL" json:"voteID"` // 投票ID
	UserID uint   `gorm:"column:user_id;type:int(11) unsigned;default:0;NOT NULL" json:"userID"` // 用户ID
	Option string `gorm:"column:option;type:char(1);NOT NULL" json:"option"`                     // 投票选项[A、B、C 、D、E、F]
}

// TableName table name
func (m *TalkRecordsVoteAnswer) TableName() string {
	return "talk_records_vote_answer"
}
