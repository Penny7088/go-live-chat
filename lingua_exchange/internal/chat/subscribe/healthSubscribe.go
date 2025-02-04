package subscribe

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/config"
	"lingua_exchange/internal/model"
)

type HealthSubscribe struct {
	config      *config.Config
	serverCache cache.ServerCache
}

func NewHealthSubscribe() *HealthSubscribe {
	return &HealthSubscribe{
		config:      config.Get(),
		serverCache: cache.NewServerCache(model.GetCacheType()),
	}
}

func (s *HealthSubscribe) Setup(ctx context.Context) error {

	log.Println("Start HealthSubscribe")

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			if err := s.serverCache.Set(ctx, s.config.App.Sid, time.Now().Unix()); err != nil {
				logger.Error(fmt.Sprintf("Websocket HealthSubscribe Report Err: %s", err.Error()))
			}
		}
	}
}
