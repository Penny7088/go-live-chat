package initial

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"golang.org/x/sync/errgroup"
	"lingua_exchange/internal/config"
	"lingua_exchange/pkg/socket"
)

type ServerOption struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type SocketServerConfig struct {
	Config   *config.Config
	ctx      context.Context
	cancel   context.CancelFunc
	engine   *gin.Engine
	server   *http.Server
	registry registry.Registry
	instance *registry.ServiceInstance
}

func (s *SocketServerConfig) String() string {
	wsAddr := ":" + strconv.Itoa(s.Config.Server.Websocket)
	return "websocket service address " + wsAddr
}

type Option func(*SocketServerConfig)

func WithRegistry(reg registry.Registry, instance *registry.ServiceInstance) Option {
	return func(s *SocketServerConfig) {
		s.registry = reg
		s.instance = instance
	}
}

func WithServerOption(opt ServerOption) Option {
	return func(s *SocketServerConfig) {
		s.server = &http.Server{
			ReadTimeout:  opt.ReadTimeout,
			WriteTimeout: opt.WriteTimeout,
			IdleTimeout:  opt.IdleTimeout,
		}
	}
}

func NewSocketServer(engine *gin.Engine, opts ...Option) *SocketServerConfig {
	ctx, cancel := context.WithCancel(context.Background())
	s := &SocketServerConfig{
		Config: config.Get(),
		ctx:    ctx,
		cancel: cancel,
		engine: engine,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *SocketServerConfig) Start() error {
	if err := s.validate(); err != nil {
		return fmt.Errorf("invalid server config: %w", err)
	}

	eg, groupCtx := errgroup.WithContext(s.ctx)
	socket.Initialize(groupCtx, eg, s.handleError)

	if err := s.registerService(); err != nil {
		return err
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Config.Server.Websocket),
		Handler: s.engine,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	log.Printf("WebSocket server starting on port: %d", s.Config.Server.Websocket)

	return s.run(eg, groupCtx, c)
}

func (s *SocketServerConfig) Stop() error {
	if s.cancel != nil {
		s.cancel()
	}

	if err := s.deregisterService(); err != nil {
		log.Printf("Failed to deregister service: %v", err)
	}

	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.server.Shutdown(ctx)
	}

	return nil
}

func (s *SocketServerConfig) validate() error {
	if s.engine == nil {
		return errors.New("gin engine is required")
	}
	if s.Config == nil {
		return errors.New("config is required")
	}
	return nil
}

func (s *SocketServerConfig) run(eg *errgroup.Group, ctx context.Context, c chan os.Signal) error {
	eg.Go(func() error {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("websocket server error: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			return s.Stop()
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	return nil
}

func (s *SocketServerConfig) registerService() error {
	if s.registry != nil && s.instance != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.registry.Register(ctx, s.instance)
	}
	return nil
}

func (s *SocketServerConfig) deregisterService() error {
	if s.registry != nil && s.instance != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.registry.Deregister(ctx, s.instance)
	}
	return nil
}

func (s *SocketServerConfig) handleError(name string) {
	if s.Config.App.Env == "prod" {
		log.Printf("WebSocket error occurred: %s", name)
	}
}
