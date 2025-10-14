package echo

import (
	"context"
	"net/http"
	"time"

	"project/internal/api"
	"project/internal/db"
	"project/internal/handler"

	"github.com/go-resty/resty/v2"
	e "github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	Address string
}

func New(log *zap.Logger) *e.Echo {
	srv := e.New()
	return srv
}

func RegisterAPIRoutes(srv *e.Echo, database *db.MemoryDB, client *resty.Client) {
	h := handler.NewApiHandler(database, client)
	api.RegisterHandlers(srv, api.NewStrictHandler(h, nil))
}

func Run(lc fx.Lifecycle, srv *e.Echo, cfg Config, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.Start(cfg.Address); err != nil && err != http.ErrServerClosed {
					log.Fatal("echo server start failed", zap.Error(err), zap.String("addr", cfg.Address))
				}
			}()
			log.Info("HTTP server starting", zap.String("addr", cfg.Address))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			log.Info("HTTP server shutting down")
			return srv.Shutdown(shutdownCtx)
		},
	})
}
