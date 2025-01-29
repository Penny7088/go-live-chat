package types

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
