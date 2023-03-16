package client

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	grpcpkg "google.golang.org/grpc"

	postv1 "sns/api/sns/post/v1"
	"sns/app/post/internal/conf"
)

// NewSnsPostServiceClient .
func NewSnsPostServiceClient(appInfo *conf.App, rr registry.Discovery, middlewares Middlewares) (postv1.PostServiceClient, error) {
	opts := []grpc.ClientOption{
		grpc.WithEndpoint("discovery:///" + appInfo.Name),
		grpc.WithDiscovery(rr),
		grpc.WithMiddleware(middlewares...),
		// grpc.WithTimeout(time.Second*2),
		grpc.WithOptions(grpcpkg.WithStatsHandler(&tracing.ClientHandler{})),
	}

	conn, err := grpc.DialInsecure(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	c := postv1.NewPostServiceClient(conn)

	return c, nil
}
