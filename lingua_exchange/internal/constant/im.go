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
