package subscribe

import (
	"context"
	"reflect"
	"sync"

	"golang.org/x/sync/errgroup"
)

var once sync.Once

// ISubscriberServer 订阅
type ISubscriberServer interface {
	Setup(ctx context.Context) error
}

// IConsume 消费
type IConsume interface {
	Call(event string, data []byte)
}

type SubscriberServers struct {
	HealthSubscribe  *HealthSubscribe  // 注册健康上报
	MessageSubscribe *MessageSubscribe // 注册消息订阅
}

func NewSubscriberServers() *SubscriberServers {
	return &SubscriberServers{
		HealthSubscribe:  NewHealthSubscribe(),
		MessageSubscribe: NewMessageSubscribe(),
	}
}

type Server struct {
	items []ISubscriberServer
}

func NewServer(servers *SubscriberServers) *Server {
	s := &Server{}

	s.binds(servers)

	return s
}

func (c *Server) binds(servers *SubscriberServers) {
	elem := reflect.ValueOf(servers).Elem()
	for i := 0; i < elem.NumField(); i++ {
		if v, ok := elem.Field(i).Interface().(ISubscriberServer); ok {
			c.items = append(c.items, v)
		}
	}
}

func (c *Server) Start(eg *errgroup.Group, ctx context.Context) {
	once.Do(func() {
		for _, process := range c.items {
			func(serv ISubscriberServer) {
				eg.Go(func() error {
					return serv.Setup(ctx)
				})
			}(process)
		}
	})
}
