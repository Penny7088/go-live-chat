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

type TalkSessionCreateRequest struct {
	TalkType   int32 `json:"talk_type"  binding:"required"`
	ReceiverID int32 `json:"receiver_id" binding:"required"`
}

type TalkSessionCreateReply struct {
	Code int                      `json:"code"` // return code
	Msg  string                   `json:"msg"`  // return information description
	Data TalkSessionCreateDetails `json:"data"`
}

type TalkSessionCreateDetails struct {
	ID         int32  `json:"id"`
	TalkType   int32  `json:"talk_type"`
	ReceiverId int32  `json:"receiver_id"`
	IsTop      int32  `json:"is_top"`
	IsDisturb  int32  `json:"is_disturb"`
	IsOnline   int32  `json:"is_online"`
	IsRobot    int32  `json:"is_robot"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	RemarkName string `json:"remark_name"`
	UnreadNum  int32  `json:"unread_num"`
	MsgText    string `json:"msg_text"`
	UpdatedAt  string `json:"updated_at"`
}

type TalkSessionDeleteRequest struct {
	SessionId int32 `json:"session_id" binding:"required"`
}

type TalkSessionDeleteReply struct {
	Result
}

type TalkSessionTopRequest struct {
	SessionId int32 `json:"session_id" binding:"required"`
	Type      int32 `json:"type" binding:"required"`
}

type TalkSessionTopReply struct {
	Result
}

type TalkSessionDisturbRequest struct {
	TalkType   int32 `json:"talk_type" binding:"required"`
	ReceiverID int32 `json:"receiver_id" binding:"required"`
	IsDisturb  int32 `json:"is_disturb" binding:"required"`
}

type TalkSessionDisturbReply struct {
	Result
}

type TalkSessionClearUnreadNumRequest struct {
	TalkType   int32 `json:"talk_type" binding:"required,oneof=1 2"`       // 对话类型
	ReceiverId int32 `json:"receiver_id" binding:"required,numeric,min=1"` // 接收者ID
}

type TalkSessionPart struct {
	ReceiverID int   `json:"receiver_id"`
	IsDisturb  int32 `json:"is_disturb"`
}
