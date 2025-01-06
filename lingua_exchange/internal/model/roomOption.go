package model

type RoomType string

const (
	RoomImGroup RoomType = "room_chat_group" // 群聊房间
	RoomExample RoomType = "room_example"    // 案例房间
)

type RoomOption struct {
	Channel  string   // 渠道分类
	RoomType RoomType // 房间类型
	Number   string   // 房间号
	Sid      string   // 网关ID
	Cid      int64    // 客户端ID
}
