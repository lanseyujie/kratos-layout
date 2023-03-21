package conf

import (
	"os"
	"runtime/debug"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is conf providers.
var ProviderSet = wire.NewSet(
	NewSourceFile,
	NewSourceKubernetes,
	NewSources,
	NewConfLoad,
	NewConfScan,
	wire.FieldsOf(new(*Bootstrap), "App", "Trace", "Server", "Data", "Registry", "Auth"),
)

const Name = "sns-post-service"

func NewSources(sf SourceFile, sk SourceKubernetes) []config.Source {
	return []config.Source{
		sf,
		sk,
	}
}

func NewConfLoad(src []config.Source) (config.Config, func(), error) {
	cfg := config.New(config.WithSource(src...))
	l := log.NewHelper(log.With(log.DefaultLogger, "module", "conf"))
	cleanup := func() {
		l.Warn("closing config resources")
		if err := cfg.Close(); err != nil {
			l.Error(err)
		}
	}

	if err := cfg.Load(); err != nil {
		return nil, cleanup, err
	}

	return cfg, cleanup, nil
}

func NewConfScan(cfg config.Config) (*Bootstrap, error) {
	bc := &Bootstrap{}
	if err := cfg.Scan(bc); err != nil {
		return nil, err
	}

	validator := func(bc *Bootstrap) error {
		if v, ok := any(bc).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}

		bc.App.Name = Name
		bc.App.Id, _ = os.Hostname()
		if bc.App.Metadata == nil {
			bc.App.Metadata = map[string]string{}
		}

		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			for _, kv := range buildInfo.Settings {
				switch kv.Key {
				case "vcs.revision":
					bc.App.Metadata["revision"] = kv.Value
				case "vcs.time":
					bc.App.Metadata["committed_at"] = kv.Value
				case "vcs.modified":
					bc.App.Metadata["dirty_build"] = kv.Value
				}
			}
		}

		return nil
	}

	if err := validator(bc); err != nil {
		return nil, err
	}

	// TODO:// 配置变更回调
	err := cfg.Watch("auth", func(key string, value config.Value) {
		log.Warnf("config changed: %s = %v\n", key, value)
	})
	if err != nil {
		return nil, err
	}

	return bc, nil
}
