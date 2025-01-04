package initial

import (
	"context"

	"golang.org/x/sync/errgroup"
	"lingua_exchange/internal/config"
	"lingua_exchange/pkg/socket"
)

type SocketServerConfig struct {
	Config *config.Config
	ctx    context.Context
}

func NewSocketServer() *SocketServerConfig {
	return &SocketServerConfig{
		Config: config.Get(),
		ctx:    context.Background(),
	}
}

func Run(socketServer *SocketServerConfig) {
	eg, groupCtx := errgroup.WithContext(socketServer.ctx)
	socket.Initialize(groupCtx, eg, func(name string) {
		if socketServer.Config.App.Env == "prod" {
			// todo  发送警告邮件
			// emailtool.SendEmail(socketServer.Config.SMTP.AdminEmail,)
		}
	})
}
