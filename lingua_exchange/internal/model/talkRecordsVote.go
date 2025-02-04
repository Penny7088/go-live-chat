package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// TalkRecordsVote 聊天对话记录（投票消息表）
type TalkRecordsVote struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	MsgID        string `gorm:"column:msg_id;type:varchar(64);NOT NULL" json:"msgID"`                                // 消息记录ID
	UserID       uint   `gorm:"column:user_id;type:int(11) unsigned;default:0;NOT NULL" json:"userID"`               // 用户ID
	Title        string `gorm:"column:title;type:varchar(64);NOT NULL" json:"title"`                                 // 投票标题
	AnswerMode   uint   `gorm:"column:answer_mode;type:tinyint(4) unsigned;default:0;NOT NULL" json:"answerMode"`    // 答题模式[0:单选;1:多选;]
	AnswerOption string `gorm:"column:answer_option;NOT NULL" json:"answerOption"`                                   // 答题选项
	AnswerNum    uint   `gorm:"column:answer_num;type:smallint(6) unsigned;default:0;NOT NULL" json:"answerNum"`     // 应答人数
	AnsweredNum  uint   `gorm:"column:answered_num;type:smallint(6) unsigned;default:0;NOT NULL" json:"answeredNum"` // 已答人数
	IsAnonymous  uint   `gorm:"column:is_anonymous;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isAnonymous"`  // 匿名投票[0:否;1:是;]
	Status       uint   `gorm:"column:status;type:tinyint(4) unsigned;default:0;NOT NULL" json:"status"`             // 投票状态[0:投票中;1:已完成;]
}

// TableName table name
func (m *TalkRecordsVote) TableName() string {
	return "talk_records_vote"
}
