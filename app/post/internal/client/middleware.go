package client

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
)

type Middlewares []middleware.Middleware

func NewClientMiddlewares(logger log.Logger) Middlewares {
	return []middleware.Middleware{
		recovery.Recovery(),
		tracing.Client(),
		logging.Client(logger),
		metadata.Client(),
		validate.Validator(),
	}
}
