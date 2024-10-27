package ws

// Message 表示客户端交互的消息体
// 包含事件名称和消息内容
type Message struct {
	Event   string `json:"event"`   // 事件名称，例如 "user_joined"
	Content any    `json:"content"` // 消息内容，可以是任意类型
}

// NewMessage 创建一个新的消息实例
// event: 事件名称
// content: 消息内容
func NewMessage(event string, content any) *Message {
	return &Message{
		Event:   event,
		Content: content,
	}
}

// SenderContent 用于推送的消息内容
type SenderContent struct {
	IsAck     bool     // 是否需要消息确认 (ACK)
	broadcast bool     // 是否是广播消息
	exclude   []int64  // 排除的用户 ID 列表（预留字段，可扩展过滤机制）
	receives  []int64  // 接收消息的用户 ID 列表
	message   *Message // 消息体，包含事件和内容
}

// NewSenderContent 创建并返回 SenderContent 的实例
func NewSenderContent() *SenderContent {
	// 预先分配切片容量，减少后续 append 调用时的内存重新分配
	return &SenderContent{
		exclude:  make([]int64, 0, 10), // 初始容量为 10
		receives: make([]int64, 0, 10),
	}
}

// SetAck 设置是否需要 ACK 确认
// value: 是否需要消息确认
// 支持链式调用
func (s *SenderContent) SetAck(value bool) *SenderContent {
	s.IsAck = value
	return s
}

// SetBroadcast 设置消息为广播类型
// value: 是否广播
// 支持链式调用
func (s *SenderContent) SetBroadcast(value bool) *SenderContent {
	s.broadcast = value
	return s
}

// SetMessage 设置消息内容
// event: 事件名称
// content: 消息内容
// 支持链式调用
func (s *SenderContent) SetMessage(event string, content any) *SenderContent {
	s.message = NewMessage(event, content)
	return s
}

// SetReceive 添加接收消息的客户端 ID 列表
// cid: 接收消息的客户端 ID
// 支持链式调用
func (s *SenderContent) SetReceive(cid ...int64) *SenderContent {
	s.receives = append(s.receives, cid...)
	return s
}

// SetExclude 设置不接收广播消息的客户端 ID 列表
// cid: 需要排除的客户端 ID
// 支持链式调用
func (s *SenderContent) SetExclude(cid ...int64) *SenderContent {
	s.exclude = append(s.exclude, cid...)
	return s
}

// IsBroadcast 判断当前消息是否是广播类型
// 返回 true 表示是广播消息，false 表示非广播消息
func (s *SenderContent) IsBroadcast() bool {
	return s.broadcast
}
