//go:generate wire .

package main

import (
	"flag"
	"os"
	"runtime/debug"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"

	"sns/app/post/internal/conf"
)

var (
	// Name is the name of the compiled software.
	Name = "sns-post"
	// Version is the version of the compiled software.
	Version = "0.0.0"
	// flagConfigPath is the config flag.
	flagConfigPath string

	id, _ = os.Hostname()
)

//nolint:gochecknoinits
func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		Version = info.Main.Version
	}
}

//nolint:gochecknoinits
func init() {
	flag.StringVar(&flagConfigPath, "config", "./configs/", "config path, eg: -config ./configs/")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	)
}

func main() {
	flag.Parse()

	// init logger
	var logger log.Logger
	{
		stdLog := log.NewStdLogger(os.Stdout)
		logger = log.With(stdLog,
			"ts", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
			"service_id", id,
			"service_name", Name,
			"service_version", Version,
			"trace_id", tracing.TraceID(),
			"span_id", tracing.SpanID(),
		)
	}

	// load config
	var bc conf.Bootstrap
	{
		cfg := config.New(config.WithSource(file.NewSource(flagConfigPath)))
		defer cfg.Close()
		if err := cfg.Load(); err != nil {
			panic(err)
		}
		if err := cfg.Scan(&bc); err != nil {
			panic(err)
		}
		if err := cfg.Watch("data.mysql", func(key string, value config.Value) {
			// TODO:// 配置变更回调
			log.Debugf("config changed: %s = %v\n", key, value)
		}); err != nil {
			log.Error(err)
		}
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err = app.Run(); err != nil {
		panic(err)
	}
}
