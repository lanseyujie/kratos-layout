package service

import (
	"github.com/google/wire"

	"sns/app/post/internal/server"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewPostService,
	// wire.Bind(new(server.RegisterServiceServer), new(*NewPostService)),
	NewServices,
)

func NewServices(post *PostService) []server.RegisterServiceServer {
	return []server.RegisterServiceServer{
		post,
	}
}
