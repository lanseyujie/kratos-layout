package server

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServerMiddlewares, NewGRPCServer, NewHTTPServer)

type RegisterServiceServer interface {
	RegisterServiceServer(srv *grpc.Server)
}

type RegisterServiceHTTPServer interface {
	RegisterServiceHTTPServer(srv *http.Server)
}
