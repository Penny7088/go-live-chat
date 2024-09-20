// Package main is the http server of the application.
package main

import (
	"github.com/zhufuyi/sponge/pkg/app"

	"lingua_exchange/cmd/lingua_exchange/initial"
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

	a := app.New(services, closes)
	a.Run()
}
