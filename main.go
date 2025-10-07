package main

import (
	"project/db"
	"project/handler"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(db.New),
		fx.Invoke(handler.RegisterHandlers),
	)
	app.Run()
}
