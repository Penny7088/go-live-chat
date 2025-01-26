package initial

import (
	"fmt"
	"log"
	"strconv"

	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/consul"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/nacos"

	"lingua_exchange/internal/config"
	routers "lingua_exchange/internal/routers"
	"lingua_exchange/internal/server"
)

func CreateServices() []app.IServer {
	var cfg = config.Get()
	var servers []app.IServer

	if httpServer, err := createHTTPServer(cfg); err != nil {
		log.Printf("Failed to create HTTP server: %v", err)
	} else {
		servers = append(servers, httpServer)
	}

	if wsServer, err := createWebSocketServer(cfg); err != nil {
		log.Printf("Failed to create WebSocket server: %v", err)
	} else {
		servers = append(servers, wsServer)
	}

	return servers
}

func createWebSocketServer(cfg *config.Config) (app.IServer, error) {
	wsRegistry, wsInstance := registerService("ws", cfg.App.Host, cfg.Server.Websocket)
	return NewSocketServer(
		routers.NewWebSocketRouter(),
		WithRegistry(wsRegistry, wsInstance),
	), nil
}

func createHTTPServer(cfg *config.Config) (app.IServer, error) {
	httpAddr := ":" + strconv.Itoa(cfg.HTTP.Port)
	httpRegistry, httpInstance := registerService("http", cfg.App.Host, cfg.HTTP.Port)
	return server.NewHTTPServer(httpAddr,
		server.WithHTTPRegistry(httpRegistry, httpInstance),
		server.WithHTTPIsProd(cfg.App.Env == "prod"),
	), nil
}

func registerService(scheme string, host string, port int) (registry.Registry, *registry.ServiceInstance) {
	var (
		instanceEndpoint = fmt.Sprintf("%s://%s:%d", scheme, host, port)
		cfg              = config.Get()
		iRegistry        registry.Registry
		instance         *registry.ServiceInstance
		err              error
		id               = cfg.App.Name + "_" + scheme + "_" + host
		logField         logger.Field
	)

	switch cfg.App.RegistryDiscoveryType {
	case "consul":
		iRegistry, instance, err = consul.NewRegistry(
			cfg.Consul.Addr,
			id,
			cfg.App.Name,
			[]string{instanceEndpoint},
		)
		if err != nil {
			panic(err)
		}
		logField = logger.Any("consulAddress", cfg.Consul.Addr)

	case "etcd":
		iRegistry, instance, err = etcd.NewRegistry(
			cfg.Etcd.Addrs,
			id,
			cfg.App.Name,
			[]string{instanceEndpoint},
		)
		if err != nil {
			panic(err)
		}
		logField = logger.Any("etcdAddress", cfg.Etcd.Addrs)

	case "nacos":
		iRegistry, instance, err = nacos.NewRegistry(
			cfg.NacosRd.IPAddr,
			cfg.NacosRd.Port,
			cfg.NacosRd.NamespaceID,
			id,
			cfg.App.Name,
			[]string{instanceEndpoint},
		)
		if err != nil {
			panic(err)
		}
		logField = logger.String("nacosAddress", fmt.Sprintf("%v:%d", cfg.NacosRd.IPAddr, cfg.NacosRd.Port))
	}

	if instance != nil {
		msg := fmt.Sprintf("register service address to %s", cfg.App.RegistryDiscoveryType)
		logger.Info(msg, logField, logger.String("id", id), logger.String("name", cfg.App.Name), logger.String("endpoint", instanceEndpoint))
		return iRegistry, instance
	}

	return nil, nil
}
