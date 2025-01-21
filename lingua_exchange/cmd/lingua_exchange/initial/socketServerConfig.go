package initial

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"lingua_exchange/internal/config"
	"lingua_exchange/pkg/socket"
)

var ErrServerClosed = errors.New("shutting down server")

type SocketServerConfig struct {
	Config *config.Config
	Ctx    context.Context
	Engine *gin.Engine
}

func NewSocketServer(engine *gin.Engine) *SocketServerConfig {
	return &SocketServerConfig{
		Config: config.Get(),
		Ctx:    context.Background(),
		Engine: engine,
	}
}

func Run(socketServer *SocketServerConfig) error {
	eg, groupCtx := errgroup.WithContext(socketServer.Ctx)
	socket.Initialize(groupCtx, eg, func(name string) {
		if socketServer.Config.App.Env == "prod" {
			// todo  发送警告邮件
			// emailtool.SendEmail(socketServer.Config.SMTP.AdminEmail,)
		}
	})

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	log.Printf("Websocket Listen Port :%d", socketServer.Config.Server.Websocket)

	return start(socketServer, eg, groupCtx, c)
}

func start(socketServer *SocketServerConfig, eg *errgroup.Group, ctx context.Context, c chan os.Signal) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", socketServer.Config.Server.Websocket),
		Handler: socketServer.Engine,
	}

	eg.Go(func() error {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	eg.Go(func() (err error) {
		defer func() {
			log.Println("Shutting down server...")

			// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
			timeCtx, timeCancel := context.WithTimeout(context.TODO(), 3*time.Second)
			defer timeCancel()

			if err := server.Shutdown(timeCtx); err != nil {
				log.Printf("Websocket Shutdown Err: %s \n", err)
			}

			err = ErrServerClosed
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			return nil
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, ErrServerClosed) {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")

	return nil
}
