package process

import (
	"context"
	"sync"
)

var once sync.Once

type IServer interface {
	Setup(ctx context.Context) error
}

type Server struct {
	items []IServer
}

// SubServers 订阅的服务列表
type SubServers struct {
	HealthSubscribe  *HealthSubscribe  // 注册健康上报
	MessageSubscribe *MessageSubscribe // 注册消息订阅
}
