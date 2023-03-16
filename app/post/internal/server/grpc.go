package server

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"sns/app/post/internal/conf"
)

// NewGRPCServer returns a gRPC server.
func NewGRPCServer(c *conf.Server, middlewares Middlewares, fns []RegisterServiceServer) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Middleware(middlewares...),
	}

	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)
	for _, fn := range fns {
		if fn != nil {
			fn.RegisterServiceServer(srv)
		}
	}

	return srv
}
