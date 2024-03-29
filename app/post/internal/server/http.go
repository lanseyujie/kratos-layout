package server

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/gorilla/handlers"

	"sns/app/post/internal/conf"
)

// NewHTTPServer returns an HTTP server.
func NewHTTPServer(c *conf.Server, middlewares Middlewares, fns []RegisterServiceHTTPServer) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(middlewares...),
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
	for _, fn := range fns {
		if fn != nil {
			fn.RegisterServiceHTTPServer(srv)
		}
	}

	srv.HandlePrefix("/q/", openapiv2.NewHandler())
	// srv.HandlePrefix("/h5/sns/post", httppkg.FileServer(httppkg.FS(os.DirFS("/srv/sns/post/dist"))))

	return srv
}
