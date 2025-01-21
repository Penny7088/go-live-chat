// Package main is the http server of the application.
package main

import (
	"log"

	"github.com/zhufuyi/sponge/pkg/app"
	"golang.org/x/sync/errgroup"
	"lingua_exchange/cmd/lingua_exchange/initial"
	routers "lingua_exchange/internal/routers"
)

// @title lingua_exchange api docs
// @description http server api docs
// @schemes http https
// @version 2.0
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type Bearer your-jwt-token to Value
func main() {
	initial.InitApp()

	services := initial.CreateServices()

	closes := initial.Close(services)

	socketServer := initial.NewSocketServer(routers.NewWebSocketRouter())

	var eg errgroup.Group
	eg.Go(func() error {
		log.Println("Starting WebSocket server...")
		return initial.Run(socketServer) // 启动 WebSocket 服务器
	})

	a := app.New(services, closes)
	a.Run()

}
