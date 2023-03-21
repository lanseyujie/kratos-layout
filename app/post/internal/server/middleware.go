package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	jwtv4 "github.com/golang-jwt/jwt/v4"

	"sns/app/post/internal/conf"
)

type WhiteList map[string]struct{}

type Middlewares []middleware.Middleware

func NewWhiteList() WhiteList {
	return WhiteList{
		// "/Login": {},
	}
}

func NewServerMiddlewares(logger log.Logger, whiteList WhiteList, authCfg *conf.Auth) Middlewares {
	return []middleware.Middleware{
		recovery.Recovery(),
		tracing.Server(),
		logging.Server(logger),
		metadata.Server(),
		validate.Validator(),
		selector.Server(
			jwt.Server(func(token *jwtv4.Token) (any, error) {
				return []byte(authCfg.Key), nil
			}, jwt.WithClaims(func() jwtv4.Claims {
				claims := jwtv4.MapClaims{}
				for k, v := range authCfg.Claims {
					claims[k] = v
				}

				return claims
			}), jwt.WithSigningMethod(jwtv4.SigningMethodHS256)),
		).Match(func(ctx context.Context, operation string) bool {
			_, exist := whiteList[operation]

			return !exist
		}).Build(),
	}
}
