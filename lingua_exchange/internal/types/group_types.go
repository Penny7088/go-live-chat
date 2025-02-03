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
	GroupID int `json:"group_id" binding:"required"`
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

// GroupRemarkUpdateRequest 表示群组备注更新请求的结构体
type GroupRemarkUpdateRequest struct {
	GroupID   int32  `json:"group_id" required:"true"`
	VisitCard string `json:"visit_card" binding:"required,max=255"`
}

type GetInviteFriendsRequest struct {
	GroupID int `json:"group_id" binding:"required"`
}

// GroupItem 表示群组项
type GroupItem struct {
	ID        int    `json:"id,omitempty"`
	GroupName string `json:"group_name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Profile   string `json:"profile,omitempty"`
	Leader    int32  `json:"leader,omitempty"`
	IsDisturb int32  `json:"is_disturb,omitempty"`
	CreatorID int32  `json:"creator_id,omitempty"`
}

// GroupListResponse 表示群组列表响应
type GroupListResponse struct {
	Items []GroupItem `json:"items,omitempty"`
}

// GroupMemberListRequest 群成员列表接口请求参数
type GroupMemberListRequest struct {
	GroupID int `json:"group_id" binding:"required"`
}

// GroupMemberListResponse 表示群组成员列表响应
type GroupMemberListResponse struct {
	Items []GroupMemberItem `json:"items"`
}

// GroupMemberItem 表示群组成员项
type GroupMemberItem struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Gender   int32  `json:"gender"`
	Leader   int32  `json:"leader"`
	IsMute   int32  `json:"is_mute,"`
	Remark   string `json:"remark"`
}

// GroupHandoverRequest 群主更换接口请求参数
type GroupHandoverRequest struct {
	GroupID int `json:"group_id" binding:"required"`
	UserID  int `json:"user_id" binding:"required"`
}

// TalkRecordExtraGroupTransfer 群主转让群消息
type TalkRecordExtraGroupTransfer struct {
	OldOwnerId   int    `json:"old_owner_id"`   // 老群主ID
	OldOwnerName string `json:"old_owner_name"` // 老群主昵称
	NewOwnerId   int    `json:"new_owner_id"`   // 新群主ID
	NewOwnerName string `json:"new_owner_name"` // 新群主昵称
}

// GroupAssignAdminRequest 分配管理员接口请求参数
type GroupAssignAdminRequest struct {
	GroupID int `json:"group_id" binding:"required"`
	UserID  int `json:"user_id" binding:"required"`
	Mode    int `json:"mode" binding:"required"`
}

// GroupNoSpeakRequest 群成员禁言接口请求参数
type GroupNoSpeakRequest struct {
	GroupID int `json:"group_id" binding:"required"`
	UserID  int `json:"user_id" binding:"required"`
	Mode    int `json:"mode" binding:"required"`
}

// GroupMuteRequest 群禁言接口请求参数
type GroupMuteRequest struct {
	GroupID int `json:"group_id" binding:"required"`
	Mode    int `json:"mode" binding:"required"` // 操作方式  1:开启全员禁言  2:解除全员禁言
}
