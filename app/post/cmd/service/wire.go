//go:build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"sns/app/post/internal/biz"
	"sns/app/post/internal/client"
	"sns/app/post/internal/conf"
	"sns/app/post/internal/data"
	"sns/app/post/internal/server"
	"sns/app/post/internal/service"
)

func newApp(cfg *conf.App, logger log.Logger, rr registry.Registrar, tp trace.TracerProvider, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.Name(cfg.Name),
		kratos.Version(cfg.Version),
		kratos.ID(cfg.Id),
		kratos.Metadata(cfg.Metadata),
		kratos.BeforeStart(func(ctx context.Context) error {
			otel.SetTracerProvider(tp)

			return nil
		}),
		kratos.Logger(logger),
		kratos.Registrar(rr),
		kratos.Server(gs),
	)
}

// wireApp init kratos application.
func wireApp(string) (*kratos.App, func(), error) {
	panic(wire.Build(
		conf.ProviderSet,
		server.ProviderSet,
		client.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}
