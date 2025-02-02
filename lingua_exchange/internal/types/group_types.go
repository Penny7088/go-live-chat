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

// GroupOutRequest 定义了退出群组请求的结构体
type GroupOutRequest struct {
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

// GroupInviteRequest 表示邀请群组请求
type GroupInviteRequest struct {
	GroupID int32  `json:"group_id" validate:"required"`
	IDs     string `json:"ids" validate:"required,ids"`
}

// GroupSettingRequest 定义了群组设置请求的结构体
type GroupSettingRequest struct {
	GroupID   int32  `json:"group_id" binding:"required"`
	GroupName string `json:"group_name" binding:"required"`
	Avatar    string `json:"avatar,omitempty"`
	Profile   string `json:"profile" binding:"max=255"`
}

// GroupRemoveMemberRequest 定义了移除群组成员请求的结构体
type GroupRemoveMemberRequest struct {
	GroupID    int32  `json:"group_id" binding:"required"`
	MembersIDs string `json:"members_ids" binding:"required,ids"`
}

type GroupDetailsRequest struct {
	GroupID    int  `json:"group_id" binding:"required"`
}

// GroupDetailResponse 群聊详情接口响应参数
type GroupDetailResponse struct {
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
	Profile   string `json:"profile"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
	IsManager bool   `json:"is_manager"`
	IsDisturb int32  `json:"is_disturb"`
	VisitCard string `json:"visit_card"`
	IsMute    int32  `json:"is_mute"`
	IsOvert   int32  `json:"is_overt"`
}
