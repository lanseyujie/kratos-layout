package client

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
)

// ProviderSet is client providers.
var ProviderSet = wire.NewSet(
	NewZapLogger,
	NewDiscoveryConsul,
	wire.Bind(new(registry.Registrar), new(*consul.Registry)),
	wire.Bind(new(registry.Discovery), new(*consul.Registry)),
	NewTracerProvider,
	NewClientMiddlewares,
	NewSnsPostServiceClient,
)
