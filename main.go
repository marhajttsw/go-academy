package main

import (
	"project/internal/db"
	appEcho "project/internal/echo"
	"project/internal/handler"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:generate go tool oapi-codegen -config ./api_gen/codegen.server.yaml ./api_gen/openapi.yaml

func main() {
	app := fx.New(
		fx.Provide(
			db.New,
			// Echo server wiring
			func() appEcho.Config { return appEcho.Config{Address: ":8080"} },
			appEcho.New,
			zap.NewExample,
		),
		fx.Invoke(
			handler.RegisterHandlers,
			appEcho.RegisterAPIRoutes,
			appEcho.Run,
		),
	)
	app.Run()
}
