package conf

import (
	"errors"
	"os"
	"runtime/debug"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is conf providers.
var ProviderSet = wire.NewSet(
	NewSourceFile,
	NewSources,
	NewConf,
	wire.FieldsOf(new(*Bootstrap), "App", "Trace", "Server", "Data", "Registry"),
)

var (
	Name    = "sns-post"
	Version = "1.0.0"
)

func NewSources(sf SourceFile) []config.Source {
	return []config.Source{
		sf,
	}
}

func NewConf(src []config.Source) (*Bootstrap, func(), error) {
	cfg := config.New(config.WithSource(src...))
	cleanup := func() {
		log.Warn("closing config resources")
		if err := cfg.Close(); err != nil {
			log.Error(err)
		}
	}

	if err := cfg.Load(); err != nil {
		return nil, cleanup, err
	}

	var bc Bootstrap
	if err := cfg.Scan(&bc); err != nil {
		return nil, cleanup, err
	}

	err := cfg.Watch("trace", func(key string, value config.Value) {
		log.Debugf("config changed: %s = %v\n", key, value)
	})
	if err != nil {
		return nil, cleanup, err
	}

	id, _ := os.Hostname()
	bc.App = &App{
		Name:     Name,
		Version:  Version,
		Id:       id,
		Metadata: map[string]string{},
	}

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		for _, kv := range buildInfo.Settings {
			bc.App.Metadata[kv.Key] = kv.Value
		}
	}

	if bc.App.Name == "" || bc.App.Version == "" || bc.App.Version == "unknown" {
		return nil, cleanup, errors.New("invalid app name or version")
	}

	return &bc, cleanup, nil
}
