package main

import (
	"project/internal/db"
	appEcho "project/internal/echo"
	"project/internal/handler"
	"project/internal/restclient"
	"time"

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
			// resty client wiring
			func() restclient.Config {
				return restclient.Config{
					BaseURL: "https://swapi.dev", // change to your external API base
					Timeout: 5 * time.Second,
					Headers: map[string]string{
						"User-Agent": "project-resty-client",
					},
				}
			},
			restclient.New,
		),
		fx.Invoke(
			handler.RegisterHandlers,
			appEcho.RegisterAPIRoutes,
			appEcho.Run,
		),
	)
	app.Run()
}
