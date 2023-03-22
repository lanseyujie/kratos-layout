package service

import (
	"github.com/google/wire"

	"sns/app/post/internal/server"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewPostService,
	NewGRPCServices,
	NewHTTPServices,
)

func NewGRPCServices(post *PostService) []server.RegisterServiceServer {
	return []server.RegisterServiceServer{
		post,
	}
}

func NewHTTPServices(post *PostService) []server.RegisterServiceHTTPServer {
	return []server.RegisterServiceHTTPServer{
		post,
	}
}
