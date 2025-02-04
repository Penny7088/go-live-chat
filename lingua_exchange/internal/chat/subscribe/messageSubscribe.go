package subscribe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc/pool"
	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/chat/consume"
	"lingua_exchange/internal/config"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/utils"
)

type MessageSubscribe struct {
	config  *config.Config
	redis   *redis.Client
	consume *consume.IMConsumer
}

func NewMessageSubscribe() *MessageSubscribe {
	return &MessageSubscribe{
		config:  config.Get(),
		redis:   model.GetRedisCli(),
		consume: consume.NewIMConsumer(),
	}
}

func (m *MessageSubscribe) SetUp(ctx context.Context) error {
	logger.Info("start subscribing message")

	go m.subscribe(ctx, []string{constant.ImTopicChat, fmt.Sprintf(constant.ImTopicChatPrivate, m.config.App.Sid)}, m.consume)
	<-ctx.Done()
	return nil
}

func (m *MessageSubscribe) subscribe(ctx context.Context, topic []string, consume IConsume) {
	sub := m.redis.Subscribe(ctx, topic...)
	defer sub.Close()

	worker := pool.New().WithMaxGoroutines(10)

	for data := range sub.Channel(redis.WithChannelHealthCheckInterval(10 * time.Second)) {
		m.handle(worker, data, consume)
	}

	worker.Wait()
}

func (m *MessageSubscribe) handle(worker *pool.Pool, data *redis.Message, consume IConsume) {
	worker.Go(func() {
		var in types.SubscribeContent
		if err := json.Unmarshal([]byte(data.Payload), &in); err != nil {
			log.Println("SubscribeContent Unmarshal Err: ", err.Error())
			return
		}

		defer func() {
			if err := recover(); err != nil {
				log.Println("MessageSubscribe Call Err: ", utils.PanicTrace(err))
			}
		}()

		consume.Call(in.Event, []byte(in.Data))
	})
}
