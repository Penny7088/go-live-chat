package types

type MemberItem struct {
	Id       string `json:"id"`
	UserId   int    `json:"user_id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Gender   int    `json:"gender"`
	Motto    string `json:"motto"`
	Leader   int    `json:"leader"`
	IsMute   int    `json:"is_mute"`
	UserCard string `json:"user_card"`
}

type PublishBaseMessageRequest struct {
	Type     string `json:"type" binding:"required"`
	Receiver struct {
		TalkType   int `json:"talk_type" binding:"required,gt=0"`   // 对话类型 1:私聊 2:群聊
		ReceiverId int `json:"receiver_id" binding:"required,gt=0"` // 好友ID或群ID
	} `json:"receiver" binding:"required"`
}

type MessageReceiver struct {
	TalkType   uint `json:"talk_type"`   // 对话类型
	ReceiverID int  `json:"receiver_id"` // 接受者ID
}

// TextMessageRequest 文本消息
type TextMessageRequest struct {
	Type     string          `json:"type"`
	Content  string          `json:"content" binding:"required"`
	QuoteID  string          `json:"quote_id"`
	Receiver MessageReceiver `json:"receiver"`
	Mentions []int32         `json:"mentions"`
}

// ImageMessageRequest 图片消息
type ImageMessageRequest struct {
	Type     string          `json:"type"`
	URL      string          `json:"url" binding:"required"`
	Width    int32           `json:"width" binding:"required"`
	Height   int32           `json:"height" binding:"required"`
	Size     int32           `json:"size" binding:"required"`
	Receiver MessageReceiver `json:"receiver"`
	QuoteID  string          `json:"quote_id"`
}

// VoiceMessageRequest 语音消息
type VoiceMessageRequest struct {
	Type     string          `json:"type"`
	URL      string          `json:"url" binding:"required"`
	Duration int32           `json:"duration" binding:"required,gt=0"`
	Size     int32           `json:"size" binding:"required"`
	Receiver MessageReceiver `json:"receiver"`
}

// VideoMessageRequest 视频文件消息
type VideoMessageRequest struct {
	Type     string          `json:"type"`
	URL      string          `json:"url" binding:"required"`
	Duration int32           `json:"duration" binding:"required,gt=0"`
	Size     int32           `json:"size" binding:"required"`
	Receiver MessageReceiver `json:"receiver"`
	Cover    string          `json:"cover"`
}

// FileMessageRequest 文件消息
type FileMessageRequest struct {
	Type     string          `json:"type"`
	UploadID string          `json:"upload_id" binding:"required"`
	Receiver MessageReceiver `json:"receiver"`
}

// CodeMessageRequest 代码消息
type CodeMessageRequest struct {
	Type     string          `json:"type"`
	Lang     string          `json:"lang" binding:"required"`
	Code     string          `json:"code" binding:"required"`
	Receiver MessageReceiver `json:"receiver"`
}

// LocationMessageRequest 位置消息
type LocationMessageRequest struct {
	Type        string          `json:"type"`
	Longitude   string          `json:"longitude" binding:"required"`
	Latitude    string          `json:"latitude" binding:"required"`
	Description string          `json:"description" binding:"required"`
	Receiver    MessageReceiver `json:"receiver"`
}

// ForwardMessageRequest 转发消息
type ForwardMessageRequest struct {
	Type       string          `json:"type"`
	Mode       int32           `json:"mode" binding:"required"`
	MessageIDs []string        `json:"message_ids" binding:"required"`
	GIDs       []int32         `json:"gids"`
	UIDs       []int32         `json:"uids"`
	Receiver   MessageReceiver `json:"receiver"`
}

// VoteMessageRequest 投票消息
type VoteMessageRequest struct {
	Type      string          `json:"type"`
	Title     string          `json:"title" binding:"required"`
	Mode      int32           `json:"mode" binding:"required"`
	Anonymous int32           `json:"anonymous" binding:"required"`
	Options   []string        `json:"options" binding:"required"`
	Receiver  MessageReceiver `json:"receiver"`
}

// LoginMessageRequest 登录消息
type LoginMessageRequest struct {
	IP       string `json:"ip"`
	Address  string `json:"address"`
	Platform string `json:"platform"`
	Agent    string `json:"agent"`
	Reason   string `json:"reason"`
}

// EmoticonMessageRequest 表情消息
type EmoticonMessageRequest struct {
	Type       string          `json:"type"`
	EmoticonID int32           `json:"emoticon_id" binding:"required"`
	Receiver   MessageReceiver `json:"receiver"`
}

// CardMessageRequest 卡片消息
type CardMessageRequest struct {
	Type     string          `json:"type"`
	UserID   int             `json:"user_id" binding:"required"`
	Receiver MessageReceiver `json:"receiver"`
}

// MixedMessageRequest 图文消息
type MixedMessageRequest struct {
	Type     string          `json:"type"`
	Items    []Item          `json:"items"`
	Receiver MessageReceiver `json:"receiver"`
	QuoteID  string          `json:"quote_id"`
}

// Item Item结构体
type Item struct {
	Type    int32  `json:"type"`
	Content string `json:"content"`
}

type AuthOption struct {
	TalkType          int
	UserId            int
	ReceiverId        uint64
	IsVerifyGroupMute bool
}
