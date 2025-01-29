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
