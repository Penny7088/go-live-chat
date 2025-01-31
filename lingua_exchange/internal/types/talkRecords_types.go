package types

// TalkRecordExtraText 文本消息
type TalkRecordExtraText struct {
	Content  string  `json:"content"`            // 文本消息
	Mentions []int32 `json:"mentions,omitempty"` // @用户ID列表
}

type Reply struct {
	UserId   int    `json:"user_id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	MsgType  int    `json:"msg_type,omitempty"` // 1:文字 2:图片
	Content  string `json:"content,omitempty"`  // 文字或图片连接
	MsgId    string `json:"msg_id,omitempty"`
}

type TalkLastMessage struct {
	MsgId      string // 消息ID
	Sequence   int    // 消息时序ID（消息排序）
	MsgType    uint   // 消息类型
	UserId     int    // 发送者ID
	ReceiverId int    // 接受者ID
	Content    string // 消息内容
	Mention    []int  // 提及列表
	CreatedAt  string // 消息发送时间
}

type TalkRecordExtraImage struct {
	Name   string `json:"name"`   // 图片名称
	Size   int    `json:"size"`   // 图片大小
	Url    string `json:"url"`    // 图片地址
	Width  int    `json:"width"`  // 图片宽度
	Height int    `json:"height"` // 图片高度
}

type TalkRecordExtraAudio struct {
	Name     string `json:"name"`     // 语音名称
	Size     int    `json:"size"`     // 语音大小
	Url      string `json:"url"`      // 语音地址
	Duration int    `json:"duration"` // 语音时长
}

type TalkRecordExtraVideo struct {
	Name     string `json:"name"`     // 视频名称
	Cover    string `json:"cover"`    // 视频封面
	Size     int    `json:"size"`     // 视频大小
	Url      string `json:"url"`      // 视频地址
	Duration int    `json:"duration"` // 视频时长
}

type TalkRecordExtraCode struct {
	Lang string `json:"lang"` // 代码语言
	Code string `json:"code"` // 代码内容
}

type TalkRecordExtraLocation struct {
	Longitude   string `json:"longitude"`   // 经度
	Latitude    string `json:"latitude"`    // 纬度
	Description string `json:"description"` // 位置描述
}

type TalkRecordExtraCard struct {
	UserId int `json:"user_id"` // 名片用户ID
}

type TalkRecordExtraMixedItem struct {
	Type    int    `json:"type"`           // 消息类型, 跟msgtype字段一致
	Content string `json:"content"`        // 消息内容。可包含图片、文字、表情等多种消息。
	Link    string `json:"link,omitempty"` // 图片跳转地址
}

type TalkRecordExtraMixed struct {
	// 消息内容。可包含图片、文字、等消息。
	Items []*TalkRecordExtraMixedItem `json:"items"` // 消息内容。可包含图片、文字、表情等多种消息。
}
