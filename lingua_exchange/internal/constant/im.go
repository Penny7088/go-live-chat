package constant

// IM 渠道分组(用于业务划分，业务间相互隔离)
const (
	// ImChannelChat 默认分组
	ImChannelChat    = "chat"    // im.Sessions.Chat.Name()
	ImChannelExample = "example" // im.Sessions.Example.Name()
)

const (
	// ImTopicChat 默认渠道消息订阅
	ImTopicChat        = "im:message:chat:all"
	ImTopicChatPrivate = "im:message:chat:%s"

	// ImTopicExample Example渠道消息订阅
	ImTopicExample        = "im:message:example:all"
	ImTopicExamplePrivate = "im:message:example:%s"
)

// 聊天模式
const (
	ChatPrivateMode = 1 // 私信模式
	ChatGroupMode   = 2 // 群聊模式
	ChatRoomMode    = 3 // 房间模式
)

// 群组
const (
	GroupMemberQuitStatusYes = 1
	GroupMemberQuitStatusNo  = 0

	GroupMemberMuteStatusYes = 1
	GroupMemberMuteStatusNo  = 0
)

// 消息类型
const (
	Text     = "text"     // 文本
	Code     = "code"     // 代码
	Location = "location" // 位置
	Emoticon = "emoticon" // 表情
	Vote     = "vote"     // 投票
	Image    = "image"    // 图片
	Voice    = "voice"    // 语音
	Video    = "video"    // 视频
	File     = "file"     // 文件
	Card     = "card"     // 卡片
	Forward  = "forward"  // 转发
	Mixed    = "mixed"    // 混合消息
)

// IM消息类型
// 1-999    自定义消息类型
// 1000-1999 系统消息类型
const (
	ChatMsgTypeText        = 1  // 文本消息
	ChatMsgTypeCode        = 2  // 代码消息
	ChatMsgTypeImage       = 3  // 图片文件
	ChatMsgTypeAudio       = 4  // 语音文件
	ChatMsgTypeVideo       = 5  // 视频文件
	ChatMsgTypeFile        = 6  // 其它文件
	ChatMsgTypeLocation    = 7  // 位置消息
	ChatMsgTypeCard        = 8  // 名片消息
	ChatMsgTypeForward     = 9  // 转发消息
	ChatMsgTypeLogin       = 10 // 登录消息
	ChatMsgTypeVote        = 11 // 投票消息
	ChatMsgTypeMixed       = 12 // 图文消息
	ChatMsgTypeGroupNotice = 13 // 撤回消息

	ChatMsgSysText                   = 1000 // 系统文本消息
	ChatMsgSysGroupCreate            = 1101 // 创建群聊消息
	ChatMsgSysGroupMemberJoin        = 1102 // 加入群聊消息
	ChatMsgSysGroupMemberQuit        = 1103 // 群成员退出群消息
	ChatMsgSysGroupMemberKicked      = 1104 // 踢出群成员消息
	ChatMsgSysGroupMessageRevoke     = 1105 // 管理员撤回成员消息
	ChatMsgSysGroupDismissed         = 1106 // 群解散
	ChatMsgSysGroupMuted             = 1107 // 群禁言
	ChatMsgSysGroupCancelMuted       = 1108 // 群解除禁言
	ChatMsgSysGroupMemberMuted       = 1109 // 群成员禁言
	ChatMsgSysGroupMemberCancelMuted = 1110 // 群成员解除禁言
	ChatMsgSysGroupNotice            = 1111 // 编辑群公告
	ChatMsgSysGroupTransfer          = 1113 // 变更群主
)

var ChatMsgTypeMapping = map[uint]string{
	ChatMsgTypeImage:                 "[图片消息]",
	ChatMsgTypeAudio:                 "[语音消息]",
	ChatMsgTypeVideo:                 "[视频消息]",
	ChatMsgTypeFile:                  "[文件消息]",
	ChatMsgTypeLocation:              "[位置消息]",
	ChatMsgTypeCard:                  "[名片消息]",
	ChatMsgTypeForward:               "[转发消息]",
	ChatMsgTypeLogin:                 "[登录消息]",
	ChatMsgTypeVote:                  "[投票消息]",
	ChatMsgTypeCode:                  "[代码消息]",
	ChatMsgTypeMixed:                 "[图文消息]",
	ChatMsgSysText:                   "[系统消息]",
	ChatMsgSysGroupCreate:            "[创建群消息]",
	ChatMsgSysGroupMemberJoin:        "[加入群消息]",
	ChatMsgSysGroupMemberQuit:        "[退出群消息]",
	ChatMsgSysGroupMemberKicked:      "[踢出群消息]",
	ChatMsgSysGroupMessageRevoke:     "[撤回消息]",
	ChatMsgSysGroupDismissed:         "[群解散消息]",
	ChatMsgSysGroupMuted:             "[群禁言消息]",
	ChatMsgSysGroupCancelMuted:       "[群解除禁言消息]",
	ChatMsgSysGroupMemberMuted:       "[群成员禁言消息]",
	ChatMsgSysGroupMemberCancelMuted: "[群成员解除禁言消息]",
	ChatMsgSysGroupNotice:            "[群公告]",
}

const (
	SubEventImMessage         = "sub.im.message"          // 对话消息通知
	SubEventImMessageKeyboard = "sub.im.message.keyboard" // 键盘输入事件通知
	SubEventImMessageRevoke   = "sub.im.message.revoke"   // 聊天消息撤销通知
	SubEventImMessageRead     = "sub.im.message.read"     // 对话消息读事件
	SubEventContactStatus     = "sub.im.contact.status"   // 用户在线状态通知
	SubEventContactApply      = "sub.im.contact.apply"    // 好友申请消息通知
	SubEventGroupJoin         = "sub.im.group.join"       // 邀请加入群聊通知
	SubEventGroupApply        = "sub.im.group.apply"      // 入群申请通知

	PushEventImMessage         = "im.message"          // 对话消息推送
	PushEventImMessageKeyboard = "im.message.keyboard" // 键盘输入事件推送
	PushEventImMessageRead     = "im.message.read"     // 对话消息读事件推送
	PushEventImMessageRevoke   = "im.message.revoke"   // 聊天消息撤销推送
	PushEventContactApply      = "im.contact.apply"    // 好友申请消息推送
	PushEventContactStatus     = "im.contact.status"   // 用户在线状态推送
	PushEventGroupApply        = "im.group.apply"      // 用户在线状态推送
)
