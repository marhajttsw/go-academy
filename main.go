package main

import (
	"net/http"
	"project/db"
	"project/handler"
	"project/httpserver"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.Provide(
			db.New,
			httpserver.NewHTTPServer,
			httpserver.NewServeMux,
			handler.NewEchoHandler,
			zap.NewExample,
		),
		fx.Invoke(
			handler.RegisterHandlers,
			func(*http.Server) {},
		),
	)
	app.Run()
}
