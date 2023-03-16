package client

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"sns/app/post/internal/conf"
)

func NewZapLogger(appInfo *conf.App) log.Logger {
	stdLog := log.NewStdLogger(os.Stdout)
	logger := log.With(stdLog,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"name", appInfo.Name,
		"version", appInfo.Version,
		"instance_id", appInfo.Id,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)

	return logger
}
