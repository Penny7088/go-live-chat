package ws

import (
	"context"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

// Session 客户端管理实例，单例模式
var Session *session

var once sync.Once

// session 渠道客户端结构
type session struct {
	Chat    *Channel // 默认分组
	Example *Channel // 示例分组

	channels map[string]*Channel
}

// Channel 获取指定名称的渠道
func (s *session) Channel(name string) (*Channel, bool) {
	val, ok := s.channels[name]
	return val, ok
}

// Initialize 初始化 Session 并启动所需的服务和协程
func Initialize(ctx context.Context, eg *errgroup.Group, fn func(name string)) {
	once.Do(func() {
		InitAck()               // 初始化 AckBuffer
		initialize(ctx, eg, fn) // 实际初始化逻辑
	})
}

// initialize 内部初始化逻辑，创建默认 Chat 和 Example 渠道并启动各项服务
func initialize(ctx context.Context, eg *errgroup.Group, fn func(name string)) {
	Session = &session{
		Chat:     NewChannel("chat", make(chan *SenderContent, 5<<20)),  // 创建 chat 渠道，缓冲区为 5MB
		Example:  NewChannel("example", make(chan *SenderContent, 100)), // 创建 example 渠道，缓冲区为 100
		channels: map[string]*Channel{},
	}

	Session.channels["chat"] = Session.Chat
	Session.channels["example"] = Session.Example

	// 延时 3 秒启动守护协程
	time.AfterFunc(3*time.Second, func() {
		// 启动健康检查协程
		eg.Go(func() error {
			defer fn("health exit")
			return health.Start(ctx)
		})

		// 启动 ack 协程
		eg.Go(func() error {
			defer fn("ack exit")
			return ack.Start(ctx)
		})

		// 启动 chat 渠道协程
		eg.Go(func() error {
			defer fn("chat exit")
			return Session.Chat.Start(ctx)
		})

		// 启动 example 渠道协程
		eg.Go(func() error {
			defer fn("example exit")
			return Session.Example.Start(ctx)
		})
	})
}
