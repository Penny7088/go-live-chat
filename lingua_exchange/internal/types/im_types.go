package types

import "lingua_exchange/internal/constant"

type SubscribeContent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type RoomOption struct {
	Channel  string            // 渠道分类
	RoomType constant.RoomType // 房间类型
	Number   string            // 房间号
	Sid      string            // 网关ID
	Cid      int64             // 客户端ID
}

type ConsumeTalkRevoke struct {
	MsgId string `json:"msg_id"`
}

type ConsumeTalkRead struct {
	SenderId   int      `json:"sender_id"`
	ReceiverId int      `json:"receiver_id"`
	MsgIds     []string `json:"msg_ids"`
}

type ConsumeTalk struct {
	TalkType   int    `json:"talk_type"`
	SenderID   int64  `json:"sender_id"`
	ReceiverID int64  `json:"receiver_id"`
	MsgId      string `json:"msg_id"`
}

type ConsumeTalkKeyboard struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
}

type ConsumeGroupJoin struct {
	Gid  int   `json:"group_id"`
	Type int   `json:"type"`
	Uids []int `json:"uids"`
}

type ConsumeGroupApply struct {
	GroupId int `json:"group_id"`
	UserId  int `json:"user_id"`
}

type ConsumeContactStatus struct {
	Status int `json:"status"`
	UserId int `json:"user_id"`
}
