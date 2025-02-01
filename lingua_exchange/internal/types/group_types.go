package types

type SearchOvertListOpt struct {
	Name   string
	UserId int
	Page   int
	Size   int
}

type GroupCreateRequest struct {
	Name   string `json:"name" binding:"required"`
	IDs    string `json:"ids" binding:"required,ids"`
	Avatar string `json:"avatar,omitempty"` // omitempty表示如果为零值则不包含在JSON中
}

// GroupCreateReply 定义了创建群组响应的结构体
type GroupCreateReply struct {
	GroupID uint64 `json:"group_id"`
}

// TalkRecordExtraGroupCreate 创建群消息
type TalkRecordExtraGroupCreate struct {
	OwnerId   uint64                        `json:"owner_id"`   // 操作人ID
	OwnerName string                        `json:"owner_name"` // 操作人昵称
	Members   []TalkRecordExtraGroupMembers `json:"members"`    // 成员列表
}

// GroupDismissRequest 定义了解散群组请求的结构体
type GroupDismissRequest struct {
	GroupID int `json:"group_id" binding:"required"` // 群组 ID，必填
}
