package model

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// TalkRecords 用户聊天记录表
type TalkRecords struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	MsgID      string    `gorm:"column:msg_id;type:varchar(64);NOT NULL" json:"msgID"`                          // 消息ID
	Sequence   int       `gorm:"column:sequence;type:int(11);NOT NULL" json:"sequence"`                         // 消息时序ID（消息排序）
	TalkType   uint      `gorm:"column:talk_type;type:int(11) unsigned;default:1;NOT NULL" json:"talkType"`     // 对话类型[1:私信;2:群聊;]
	MsgType    uint      `gorm:"column:msg_type;type:int(11) unsigned;default:1;NOT NULL" json:"msgType"`       // 消息类型[1:文本消息;2:文件消息;3:会话消息;4:代码消息;5:投票消息;6:群公告;7:好友申请;8:登录通知;9:入群消息/退群消息;]
	UserID     uint64    `gorm:"column:user_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"userID"`      // 发送者ID（0:代表系统消息 >0: 用户ID）
	ReceiverID uint      `gorm:"column:receiver_id;type:int(11) unsigned;default:0;NOT NULL" json:"receiverID"` // 接收者ID（用户ID 或 群ID）
	IsRevoke   uint      `gorm:"column:is_revoke;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isRevoke"`  // 是否撤回[0:否;1:是;]
	IsMark     uint      `gorm:"column:is_mark;type:tinyint(4) unsigned;default:0;NOT NULL" json:"isMark"`      // 是否重要[0:否;1:是;]
	QuoteID    string    `gorm:"column:quote_id;type:varchar(64);NOT NULL" json:"quoteID"`                      // 引用消息ID
	Extra      string    `gorm:"column:extra;default:{}" json:"extra"`                                          // 消息扩展字段
	CreatedAt  time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`                                  // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`                                  // 更新时间
}
