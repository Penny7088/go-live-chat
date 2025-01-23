package types

type TalkSessionItem struct {
	ID         int32  `json:"id"`
	TalkType   int32  `json:"talk_type"`
	ReceiverID int32  `json:"receiver_id"`
	IsTop      int32  `json:"is_top"`
	IsDisturb  int32  `json:"is_disturb"`
	IsOnline   int32  `json:"is_online"`
	IsRobot    int32  `json:"is_robot"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Remark     string `json:"remark"`
	UnreadNum  int32  `json:"unread_num"`
	MsgText    string `json:"msg_text"`
	UpdatedAt  string `json:"updated_at"`
}

type TalkSessionItemsReply struct {
	Code int               `json:"code"` // return code
	Msg  string            `json:"msg"`  // return information description
	Data []TalkSessionItem `json:"data"`
}
