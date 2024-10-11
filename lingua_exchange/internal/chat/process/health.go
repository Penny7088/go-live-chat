package process

import (
	"context"
	"fmt"
	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/config"
	"time"
)

type HealthSubscribe struct {
	storage *cache.WsServerIDCache
}

func NewHealthSubscribe(cache *cache.WsServerIDCache) *HealthSubscribe {
	return &HealthSubscribe{storage: cache}
}

func (s *HealthSubscribe) Setup(ctx context.Context) error {

	logger.Info("Start HealthSubscribe ... ")
	config := config.Get()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			if err := s.storage.Set(ctx, config.ServerId(), time.Now().Unix()); err != nil {
				logger.Error(fmt.Sprintf("Websocket HealthSubscribe Report Err: %s", err.Error()))
			}
		}
	}
}
