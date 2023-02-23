package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/gorilla/handlers"

	postv1 "sns/api/sns/post/v1"
	"sns/app/post/internal/conf"
	"sns/app/post/internal/service"
)

// NewHTTPServer returns an HTTP server.
func NewHTTPServer(c *conf.Server, post *service.PostService, logger log.Logger) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			// tracing.Server(),
			logging.Server(logger),
			metadata.Server(),
			validate.Validator(),
		),
		http.Filter(
			handlers.CORS(
				handlers.AllowedHeaders([]string{"*"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
				handlers.AllowedOrigins([]string{"*"}),
				handlers.AllowCredentials(),
			),
		),
	}

	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)
	postv1.RegisterPostServiceHTTPServer(srv, post)

	srv.HandlePrefix("/q/", openapiv2.NewHandler())
	// srv.HandlePrefix("/h5/sns/post", httppkg.FileServer(httppkg.FS(os.DirFS("/srv/sns/post/dist"))))

	return srv
}
